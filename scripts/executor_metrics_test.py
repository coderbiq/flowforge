#!/usr/bin/env python3
"""Tests for scripts/executor_metrics.py (ticket 01-metrics-extractor).

Metric-level tests replay known shapes against an in-memory fixture DB
(sqlite3 :memory:). CLI-level tests use a temp-file fixture DB because the
script opens database files read-only via a file:...?mode=ro URI.

Run: python3 -m unittest discover -s scripts -p 'executor_metrics_test.py' -v
"""

import contextlib
import io
import json
import os
import re
import sqlite3
import tempfile
import unittest

import executor_metrics as em

BASE_TS = 1789200000000  # fixture epoch-ms base (2026-09-14T... local-independent)

SESSION_COLS = ("id, agent, model, directory, tokens_input, tokens_output, "
                "tokens_cache_read, tokens_cache_write, cost, time_created, time_updated")
FLASH_MODEL = '{"id":"gemini-3.8-flash-high","providerID":"erasebg-gemini","variant":"default"}'
FLAGSHIP_MODEL = '{"id":"glm-4.7","providerID":"zhipu","variant":"default"}'


SCHEMA = (
    "CREATE TABLE session (id TEXT PRIMARY KEY, agent TEXT, model TEXT, "
    "directory TEXT, tokens_input INTEGER, tokens_output INTEGER, "
    "tokens_cache_read INTEGER, tokens_cache_write INTEGER, cost REAL, "
    "time_created INTEGER, time_updated INTEGER)",
    "CREATE TABLE part (id TEXT PRIMARY KEY, session_id TEXT, "
    "time_created INTEGER, data TEXT)",
)


def make_conn():
    conn = sqlite3.connect(":memory:")
    for ddl in SCHEMA:
        conn.execute(ddl)
    return conn


def add_session(conn, sid, tc=BASE_TS, dur_min=10.0, model=FLASH_MODEL,
                directory="/work/tangram-v2", agent="flowforge-implementer",
                ti=1000, to=2000, tcr=5000):
    conn.execute(f"INSERT INTO session ({SESSION_COLS}) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
                 (sid, agent, model, directory, ti, to, tcr, 0, 0.0,
                  tc, tc + int(dur_min * 60000)))


_part_seq = 0


def add_part(conn, sid, payload, tc=None):
    global _part_seq
    _part_seq += 1
    if tc is None:
        tc = BASE_TS + _part_seq
    if isinstance(payload, (dict, list)):
        payload = json.dumps(payload)
    conn.execute("INSERT INTO part (id, session_id, time_created, data) VALUES (?,?,?,?)",
                 (f"prt_{_part_seq:08d}", sid, tc, payload))
    return tc


def text_part(txt):
    return {"type": "text", "text": txt}


def step_part():
    return {"type": "step-start"}


def tool_part(tool="bash", command=None, output=None):
    state = {"status": "completed"}
    if command is not None:
        state["input"] = {"command": command}
    if output is not None:
        state["output"] = output
    return {"type": "tool", "tool": tool, "state": state}


def dispatch_and_report_parts(ticket_line="ticket docs/proposals/timeline/issues/06-x.md here",
                              report="STATUS: COMPLETED\n\n### Summary\ndone"):
    return [text_part(f"你是 implementer。{ticket_line}"), text_part(report)]


def run_cli(argv):
    out, err = io.StringIO(), io.StringIO()
    with contextlib.redirect_stdout(out), contextlib.redirect_stderr(err):
        code = em.main(argv)
    return code, out.getvalue(), err.getvalue()


def write_db_file(path, populate):
    if os.path.exists(path):
        os.remove(path)
    conn = sqlite3.connect(path)
    for ddl in SCHEMA:
        conn.execute(ddl)
    populate(conn)
    conn.commit()
    conn.close()


class VerdictTests(unittest.TestCase):
    def test_four_way_classification(self):
        self.assertEqual(em.classify_verdict("...\nBUILD SUCCESSFUL in 3s\n"), "ok")
        self.assertEqual(em.classify_verdict("FAILURE: Build failed\nBUILD FAILED\n"), "fail")
        self.assertEqual(em.classify_verdict("Command failed with exit code 1"), "fail")
        self.assertEqual(em.classify_verdict("Command failed with exit code 127: x"), "fail")
        self.assertEqual(em.classify_verdict("Command failed with exit code 0"), "unknown")
        self.assertEqual(em.classify_verdict("ls: 3 files total, no marker"), "unknown")
        self.assertEqual(em.classify_verdict(None), "unknown")

    def test_build_successful_wins_over_exit_code(self):
        self.assertEqual(em.classify_verdict("BUILD SUCCESSFUL\n(exit code 0)"), "ok")


class NormalizeTests(unittest.TestCase):
    def test_whitespace_collapsed_and_trimmed(self):
        self.assertEqual(em.normalize_command("  go   test\n\t./...  "), "go test ./...")

    def test_truncated_to_120_chars(self):
        self.assertEqual(len(em.normalize_command("x" * 500)), 120)
        # identical first 120 chars => same normalized key even if tails differ
        self.assertEqual(em.normalize_command("a" * 120 + "XY"),
                         em.normalize_command("a" * 120 + "ZW"))


class SessionMetricsTests(unittest.TestCase):
    def one(self, conn, **kw):
        rows, corrupt = em.collect_metrics(conn, kw.pop("project", "tangram-v2"), **kw)
        self.assertEqual(corrupt, 0)
        self.assertEqual(len(rows), 1)
        return rows[0]

    def test_counts_and_token_columns(self):
        conn = make_conn()
        add_session(conn, "ses_a", ti=1100, to=2200, tcr=3300, dur_min=1.5)
        for p in dispatch_and_report_parts():
            add_part(conn, "ses_a", p)
        add_part(conn, "ses_a", step_part(), tc=BASE_TS + 10)
        add_part(conn, "ses_a", step_part(), tc=BASE_TS + 11)
        add_part(conn, "ses_a", step_part(), tc=BASE_TS + 12)
        add_part(conn, "ses_a", tool_part("bash", "go build ./...", "ok"), tc=BASE_TS + 20)
        add_part(conn, "ses_a", tool_part("edit"), tc=BASE_TS + 21)
        add_part(conn, "ses_a", tool_part("bash", "ls -la", "3 files"), tc=BASE_TS + 22)
        m = self.one(conn)
        self.assertEqual(m["session"], "ses_a")
        self.assertEqual(m["steps"], 3)
        self.assertEqual(m["tools"], 3)
        self.assertEqual(m["bash_n"], 2)
        self.assertEqual(m["in_tok"], 1100)
        self.assertEqual(m["out_tok"], 2200)
        self.assertEqual(m["cacheR_tok"], 3300)
        self.assertEqual(m["dur_min"], 1.5)
        self.assertEqual(m["model"], "gemini-3.8-flash-high")
        self.assertEqual(m["provider"], "erasebg-gemini")
        self.assertEqual(m["ticket"], "proposals/timeline/issues/06-x.md")
        self.assertEqual(m["status"], "COMPLETED")
        self.assertRegex(m["ts"], r"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[+-]\d{2}:\d{2}$")

    def test_status_missing_records_none(self):
        conn = make_conn()
        add_session(conn, "ses_b")
        add_part(conn, "ses_b", text_part("STATUS: BLOCKED midway context"), tc=BASE_TS + 1)
        add_part(conn, "ses_b", text_part("final text without any marker"), tc=BASE_TS + 2)
        self.assertEqual(self.one(conn)["status"], "none")

    def test_status_from_last_text_part_first_match(self):
        conn = make_conn()
        add_session(conn, "ses_c")
        add_part(conn, "ses_c", text_part("earlier STATUS: BLOCKED"), tc=BASE_TS + 1)
        add_part(conn, "ses_c", text_part("STATUS: COMPLETED ... later quotes STATUS: INCONCLUSIVE"), tc=BASE_TS + 2)
        self.assertEqual(self.one(conn)["status"], "COMPLETED")

    def test_rep_max_counts_normalized_commands(self):
        conn = make_conn()
        add_session(conn, "ses_d")
        for p in dispatch_and_report_parts():
            add_part(conn, "ses_d", p)
        add_part(conn, "ses_d", tool_part("bash", "go test ./...", "exit code 1"), tc=BASE_TS + 1)
        add_part(conn, "ses_d", tool_part("bash", "  go  test\n./... ", "exit code 2"), tc=BASE_TS + 2)
        add_part(conn, "ses_d", tool_part("bash", "git status", "clean"), tc=BASE_TS + 3)
        m = self.one(conn)
        self.assertEqual(m["rep_max"], 2)

    def test_fail_streak_unknown_skipped_ok_resets(self):
        conn = make_conn()
        add_session(conn, "ses_e")
        for p in dispatch_and_report_parts():
            add_part(conn, "ses_e", p)
        seq = [
            ("c1", "Command failed with exit code 1"),   # fail -> 1
            ("c2", "Command failed with exit code 2"),   # fail -> 2
            ("c3", "no recognizable marker"),            # unknown -> skipped
            ("c4", "Command failed with exit code 3"),   # fail -> 3
            ("c5", "BUILD SUCCESSFUL"),                  # ok -> reset
            ("c6", "BUILD FAILED"),                      # fail -> 1
        ]
        for i, (cmd, out) in enumerate(seq):
            add_part(conn, "ses_e", tool_part("bash", cmd, out), tc=BASE_TS + i + 1)
        m = self.one(conn)
        self.assertEqual(m["fail_streak"], 3)

    def test_corrupt_part_json_skipped_and_counted(self):
        conn = make_conn()
        add_session(conn, "ses_f")
        add_part(conn, "ses_f", text_part("docs/proposals/p/issues/09-y.md"))
        add_part(conn, "ses_f", step_part(), tc=BASE_TS + 1)
        add_part(conn, "ses_f", '{"type": "tool", "tool": "bash", "state": {"in', tc=BASE_TS + 2)
        add_part(conn, "ses_f", tool_part("bash", "true", "exit code 0"), tc=BASE_TS + 3)
        rows, corrupt = em.collect_metrics(conn, "tangram-v2")
        self.assertEqual(corrupt, 1)
        self.assertEqual(len(rows), 1)
        self.assertEqual(rows[0]["steps"], 1)
        self.assertEqual(rows[0]["tools"], 1)

    def test_filter_agent_directory_and_order(self):
        conn = make_conn()
        add_session(conn, "ses_g1", tc=BASE_TS + 100)
        add_session(conn, "ses_g2", tc=BASE_TS + 50)
        add_session(conn, "ses_g3", tc=BASE_TS + 10, agent="flowforge-reviewer")
        add_session(conn, "ses_g4", tc=BASE_TS + 10, directory="/work/other-proj")
        add_session(conn, "ses_g5", tc=BASE_TS + 10, directory="/work/tangram-v2-old")
        for sid in ("ses_g1", "ses_g2", "ses_g3", "ses_g4", "ses_g5"):
            add_part(conn, sid, dispatch_and_report_parts()[1])
        rows, _ = em.collect_metrics(conn, "tangram-v2")
        self.assertEqual([r["session"] for r in rows], ["ses_g2", "ses_g1"])

    def test_epoch_default_dash(self):
        conn = make_conn()
        add_session(conn, "ses_h")
        add_part(conn, "ses_h", dispatch_and_report_parts()[1])
        self.assertEqual(self.one(conn)["epoch"], "-")


class EpochTests(unittest.TestCase):
    def test_parse_and_label(self):
        epochs = em.parse_epochs("pre:<100,mid:100-200,post:>200")
        cases = {50: "pre", 99: "pre", 100: "mid", 150: "mid", 200: "mid", 201: "post", 300: "post"}
        for ts, want in cases.items():
            self.assertEqual(em.epoch_label(ts, epochs), want, f"ts={ts}")
        self.assertIsNone(em.parse_epochs(None))

    def test_bad_spec_raises(self):
        with self.assertRaises(em.MetricsError):
            em.parse_epochs("bad")


class CostTests(unittest.TestCase):
    def row_for(self, model, ti, to, tcr, prices=None):
        conn = make_conn()
        add_session(conn, "ses_cost", model=model, ti=ti, to=to, tcr=tcr)
        add_part(conn, "ses_cost", dispatch_and_report_parts()[1])
        rows, _ = em.collect_metrics(conn, "tangram-v2", prices=prices)
        return rows[0]

    def test_flash_tier_proxy_price(self):
        # 1M in + 1M out + 1M cacheR at flash prices 0.30/2.50/0.075 -> 2.875
        m = self.row_for(FLASH_MODEL, 1_000_000, 1_000_000, 1_000_000)
        self.assertEqual(m["est_cost"], 2.875)

    def test_flagship_tier_proxy_price(self):
        m = self.row_for(FLAGSHIP_MODEL, 1_000_000, 1_000_000, 1_000_000)
        self.assertEqual(m["est_cost"], round(0.60 + 2.20 + 0.113, 4))

    def test_price_override(self):
        m = self.row_for(FLASH_MODEL, 1_000_000, 1_000_000, 1_000_000,
                         prices={"flash": (1.0, 1.0, 1.0)})
        self.assertEqual(m["est_cost"], 3.0)


class CliExtractTests(unittest.TestCase):
    def setUp(self):
        self._tmp = tempfile.TemporaryDirectory()
        self.dir = self._tmp.name
        self.db = os.path.join(self.dir, "fixture.db")
        self.out = os.path.join(self.dir, "obs.md")
        self.addCleanup(self._tmp.cleanup)

    def populate(self, conn):
        add_session(conn, "ses_x1", tc=BASE_TS + 100, ti=100, to=200, tcr=400)
        for p in dispatch_and_report_parts():
            add_part(conn, "ses_x1", p)
        add_part(conn, "ses_x1", step_part(), tc=BASE_TS + 101)
        add_part(conn, "ses_x1", tool_part("bash", "go build ./...", "BUILD SUCCESSFUL"), tc=BASE_TS + 102)
        add_session(conn, "ses_x2", tc=BASE_TS + 200)
        for p in dispatch_and_report_parts(report="STATUS: BLOCKED\n\nblocker: x"):
            add_part(conn, "ses_x2", p)
        add_session(conn, "ses_x3", tc=BASE_TS + 300)

    def test_writes_header_and_rows(self):
        write_db_file(self.db, self.populate)
        code, out, err = run_cli(["extract", "--db", self.db, "--project", "tangram-v2",
                                  "--out", self.out])
        self.assertEqual(code, 0, err)
        with open(self.out, encoding="utf-8") as f:
            content = f.read()
        self.assertTrue(content.endswith("\n"))
        lines = content.splitlines()
        self.assertIn("| session | ts | dur_min | model | provider | ticket | epoch | "
                      "steps | tools | bash_n | rep_max | fail_streak | in_tok | out_tok | "
                      "cacheR_tok | est_cost | status |", lines)
        rows = [ln for ln in lines if ln.startswith("| ses_")]
        self.assertEqual(len(rows), 3)
        row1 = next(ln for ln in rows if "ses_x1" in ln)
        self.assertIn("| 1 | 1 | 1 | 1 | 0 |", row1)  # steps tools bash_n rep_max fail_streak
        self.assertIn("COMPLETED", row1)
        self.assertIn("erasebg-gemini", row1)
        row2 = next(ln for ln in rows if "ses_x2" in ln)
        self.assertIn("BLOCKED", row2)

    def test_idempotent_rerun_and_incremental_append(self):
        write_db_file(self.db, self.populate)
        args = ["extract", "--db", self.db, "--project", "tangram-v2", "--out", self.out]
        self.assertEqual(run_cli(args)[0], 0)
        with open(self.out, encoding="utf-8") as f:
            first = f.read()
        code, out, err = run_cli(args)
        self.assertEqual(code, 0, err)
        with open(self.out, encoding="utf-8") as f:
            second = f.read()
        self.assertEqual(first, second, "rerun must not change observations.md")

        def add_later(conn):
            self.populate(conn)
            add_session(conn, "ses_x4", tc=BASE_TS + 400)
            add_part(conn, "ses_x4", dispatch_and_report_parts()[1])

        write_db_file(self.db, add_later)
        code, out, err = run_cli(args)
        self.assertEqual(code, 0, err)
        with open(self.out, encoding="utf-8") as f:
            third = f.read()
        self.assertEqual(len([ln for ln in third.splitlines() if ln.startswith("| ses_")]), 4)
        self.assertTrue(third.startswith(first), "append-only: previous bytes unchanged")

    def test_since_inclusive_boundary(self):
        write_db_file(self.db, self.populate)
        code, _, err = run_cli(["extract", "--db", self.db, "--project", "tangram-v2",
                                "--out", self.out, "--since", str(BASE_TS + 200)])
        self.assertEqual(code, 0, err)
        with open(self.out, encoding="utf-8") as f:
            rows = [ln for ln in f.read().splitlines() if ln.startswith("| ses_")]
        self.assertEqual(len(rows), 2)  # ses_x2 (== since, inclusive) and ses_x3
        self.assertTrue(any("ses_x2" in ln for ln in rows))

    def test_corrupt_json_warns_on_stderr(self):
        def populate_with_corruption(conn):
            self.populate(conn)
            add_part(conn, "ses_x1", '{"type": "tool", "broken"', tc=BASE_TS + 105)
        write_db_file(self.db, populate_with_corruption)
        code, _, err = run_cli(["extract", "--db", self.db, "--project", "tangram-v2",
                                "--out", self.out])
        self.assertEqual(code, 0)
        self.assertIn("warn: skipped 1", err)

    def test_price_override_and_epochs_flags(self):
        write_db_file(self.db, self.populate)
        code, _, err = run_cli([
            "extract", "--db", self.db, "--project", "tangram-v2", "--out", self.out,
            "--epochs", f"a:<{BASE_TS + 150},b:{BASE_TS + 150}-{BASE_TS + 250},c:>{BASE_TS + 250}",
            "--price-override", '{"flash": [1, 1, 1]}'])
        self.assertEqual(code, 0, err)
        with open(self.out, encoding="utf-8") as f:
            rows = [ln for ln in f.read().splitlines() if ln.startswith("| ses_")]
        self.assertIn("| a |", next(ln for ln in rows if "ses_x1" in ln))
        self.assertIn("| b |", next(ln for ln in rows if "ses_x2" in ln))
        self.assertIn("| c |", next(ln for ln in rows if "ses_x3" in ln))


class CliFailureTests(unittest.TestCase):
    def setUp(self):
        self._tmp = tempfile.TemporaryDirectory()
        self.dir = self._tmp.name
        self.out = os.path.join(self.dir, "obs.md")
        self.addCleanup(self._tmp.cleanup)

    def test_missing_db_errors_without_creating_out(self):
        code, _, err = run_cli(["extract", "--db", os.path.join(self.dir, "nope.db"),
                                "--project", "tangram-v2", "--out", self.out])
        self.assertEqual(code, 1)
        self.assertIn("error:", err)
        self.assertFalse(os.path.exists(self.out))

    def test_non_sqlite_file_errors_without_creating_out(self):
        junk = os.path.join(self.dir, "junk.db")
        with open(junk, "wb") as f:
            f.write(b"this is definitely not a sqlite database" * 10)
        code, _, err = run_cli(["extract", "--db", junk, "--project", "tangram-v2",
                                "--out", self.out])
        self.assertEqual(code, 1)
        self.assertIn("error:", err)
        self.assertFalse(os.path.exists(self.out))

    def test_sqlite_without_opencode_tables_errors(self):
        empty = os.path.join(self.dir, "empty.db")
        conn = sqlite3.connect(empty)
        conn.execute("CREATE TABLE other (x)")
        conn.commit()
        conn.close()
        code, _, err = run_cli(["extract", "--db", empty, "--project", "tangram-v2",
                                "--out", self.out])
        self.assertEqual(code, 1)
        self.assertIn("error:", err)
        self.assertFalse(os.path.exists(self.out))

    def test_out_is_directory_errors_cleanly(self):
        db = os.path.join(self.dir, "fixture.db")
        write_db_file(db, lambda conn: None)
        code, _, err = run_cli(["extract", "--db", db, "--project", "tangram-v2",
                                "--out", self.dir])
        self.assertEqual(code, 1)
        self.assertIn("cannot write observations file", err)


if __name__ == "__main__":
    unittest.main()
