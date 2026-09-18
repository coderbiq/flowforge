#!/usr/bin/env python3
"""Tests for scripts/executor_metrics.py (tickets 01/02/03).

Metric-level tests replay known shapes against an in-memory fixture DB
(sqlite3 :memory:). CLI-level tests use a temp-file fixture DB because the
script opens database files read-only via a file:...?mode=ro URI.
Gate tests (ticket 03) drive the public gate/evaluation functions, the
report renderer and the DECISION.md template artifact.

Run: python3 -m unittest discover -s scripts -p 'executor_metrics_test.py' -v
"""

import contextlib
import io
import json
import os
import pathlib
import re
import sqlite3
import tempfile
import unittest
from datetime import datetime, timezone

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
        self.assertIn("| session | agent | ts | dur_min | model | provider | ticket | "
                      "epoch | steps | tools | bash_n | rep_max | fail_streak | in_tok | "
                      "out_tok | cacheR_tok | est_cost | status |", lines)
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


class ObsParsingTests(unittest.TestCase):
    """report: parse observations.md table rows (header-driven)."""

    def obs_text(self, *rows):
        lines = list(em.HEADER_LINES)
        for r in rows:
            lines.append(em.format_row(r))
        return "\n".join(lines) + "\n"

    def row(self, **kw):
        base = {"session": "ses_r1", "ts": "2026-09-13T12:35:44+08:00",
                "dur_min": 1.5, "model": "gemini-3.8-flash-high",
                "provider": "erasebg-gemini", "ticket": "proposals/p/issues/01-x.md",
                "epoch": "-", "steps": 10, "tools": 8, "bash_n": 4, "rep_max": 1,
                "fail_streak": 0, "in_tok": 1000, "out_tok": 2000,
                "cacheR_tok": 5000, "est_cost": 0.01, "status": "COMPLETED"}
        base.update(kw)
        return base

    def test_rows_parsed_header_driven_with_types(self):
        rows = em.parse_obs_rows(self.obs_text(self.row(), self.row(session="ses_r2")))
        self.assertEqual(len(rows), 2)
        r = rows[0]
        self.assertEqual(r["session"], "ses_r1")
        self.assertEqual(r["steps"], 10)          # int
        self.assertEqual(r["cacheR_tok"], 5000)   # int
        self.assertEqual(r["dur_min"], 1.5)       # float
        self.assertEqual(r["ticket"], "proposals/p/issues/01-x.md")

    def test_no_table_header_is_error(self):
        with self.assertRaises(em.MetricsError):
            em.parse_obs_rows("# only prose\n\nno table here\n")

    def test_column_count_mismatch_is_error(self):
        text = self.obs_text(self.row()) + "| ses_bad | only | three |\n"
        with self.assertRaises(em.MetricsError):
            em.parse_obs_rows(text)

    def test_ts_to_ms_roundtrip_and_bad_value(self):
        ms = em.ts_ms("2026-09-13T12:35:44+08:00")
        self.assertEqual(ms, 1789274144000)
        with self.assertRaises(em.MetricsError):
            em.ts_ms("not-a-timestamp")


class ReportEpochTests(unittest.TestCase):
    """report: epoch placement of observation rows; uncovered -> unclassified."""

    def test_boundary_ts_goes_to_right_interval(self):
        # "<1000" strict, "1000-2000" inclusive both ends, ">2000" excludes 2000.
        epochs = em.parse_epochs("pre:<1000,mid:1000-2000,post:>2000")
        cases = {999: "pre", 1000: "mid", 1500: "mid", 2000: "mid", 2001: "post"}
        for ts, want in cases.items():
            self.assertEqual(em.report_epoch(ts, epochs), want, f"ts={ts}")

    def test_uncovered_and_no_epochs_are_unclassified(self):
        epochs = em.parse_epochs("mid:1000-2000")
        self.assertEqual(em.report_epoch(500, epochs), em.UNCLASSIFIED)
        self.assertEqual(em.report_epoch(3000, epochs), em.UNCLASSIFIED)
        self.assertEqual(em.report_epoch(1500, None), em.UNCLASSIFIED)

    def test_real_baseline_boundaries(self):
        epochs = em.parse_epochs(
            "pre-hardening:<1789291500000,"
            "provider-switch:1789291500000-1789298700000,hardened:>1789298700000")
        incident = em.ts_ms("2026-09-13T12:35:44+08:00")
        provider_first = em.ts_ms("2026-09-13T17:38:09+08:00")
        boundary_1725 = em.ts_ms("2026-09-13T17:25:00+08:00")
        boundary_1925 = em.ts_ms("2026-09-13T19:25:00+08:00")
        self.assertEqual(em.report_epoch(incident, epochs), "pre-hardening")
        self.assertEqual(em.report_epoch(provider_first, epochs), "provider-switch")
        # "<T" excludes T; "T1-T2" is right-closed (owns its endpoints);
        # ">T" excludes T — a ts exactly at 19:25:00 belongs to provider-switch.
        self.assertEqual(em.report_epoch(boundary_1725, epochs), "provider-switch")
        self.assertEqual(em.report_epoch(boundary_1925, epochs), "provider-switch")
        self.assertEqual(
            em.report_epoch(em.ts_ms("2026-09-13T19:25:00.001+08:00"), epochs),
            "hardened")


TICKET_S = """# 01: small
## Changes
- [x] 1. one change
- [x] 2. two change
## Completion evidence
- cmd: `./gradlew :app:test --tests "*Foo*"`
  - output: BUILD SUCCESSFUL
"""

TICKET_M_BY_COUNT = """# 02: three changes
## Changes
- [x] 1. one
- [ ] 2. two
- [x] 3. three
## Completion evidence
- cmd: `./gradlew :app:test --tests "*Foo*"`
"""

TICKET_L_BY_COUNT = """# 03: four changes
## Changes
- [x] 1. one
- [x] 2. two
- [x] 3. three
- [ ] 4. four
## Completion evidence
- cmd: `./gradlew :app:test --tests "*Foo*"`
"""

TICKET_L_BY_INTEGRATION = """# 04: heavy cmd, few changes
## Changes
- [x] 1. one
## Completion evidence
- cmd: `./gradlew :app:integrationTest`
"""

TICKET_L_BY_RERUN = """# 05: rerun
## Changes
- [x] 1. one
- [x] 2. two
## Completion evidence
- cmd: `./gradlew :app:test --tests "*Foo*" --rerun`
"""

TICKET_M_BY_NONUNIT_CMD = """# 06: build-only evidence
## Changes
- [x] 1. one
- [x] 2. two
## Completion evidence
- cmd: `./gradlew :app:compileKotlin`
"""


class StrataTests(unittest.TestCase):
    """report: objective S/M/L ticket-difficulty proxy from ticket files."""

    def test_s_boundary_two_changes_unit_cmd(self):
        self.assertEqual(em.classify_stratum_text(TICKET_S), "S")
        # 2 changes with pnpm test / go test are unit-level too
        self.assertEqual(em.classify_stratum_text(
            TICKET_S.replace("`./gradlew :app:test --tests \"*Foo*\"`", "`pnpm test`")), "S")
        self.assertEqual(em.classify_stratum_text(
            TICKET_S.replace("`./gradlew :app:test --tests \"*Foo*\"`", "`go test ./pkg/...`")), "S")

    def test_l_boundary_four_changes(self):
        self.assertEqual(em.classify_stratum_text(TICKET_L_BY_COUNT), "L")

    def test_l_by_heavy_cmd_regardless_of_count(self):
        self.assertEqual(em.classify_stratum_text(TICKET_L_BY_INTEGRATION), "L")
        self.assertEqual(em.classify_stratum_text(TICKET_L_BY_RERUN), "L")
        self.assertIn("E2E", em.HEAVY_CMD_RE.pattern)

    def test_m_for_rest(self):
        self.assertEqual(em.classify_stratum_text(TICKET_M_BY_COUNT), "M")
        self.assertEqual(em.classify_stratum_text(TICKET_M_BY_NONUNIT_CMD), "M")
        # no command-shaped code spans at all -> cannot confirm unit-level -> M
        self.assertEqual(em.classify_stratum_text(
            "# t\n## Changes\n- [x] 1. one\n- [x] 2. two\nno cmds here\n"), "M")


class ResolveStrataTests(unittest.TestCase):
    """report: map observations ticket paths to files under --project-root."""

    def setUp(self):
        self._tmp = tempfile.TemporaryDirectory()
        self.root = pathlib.Path(self._tmp.name)
        self.addCleanup(self._tmp.cleanup)

    def write_ticket(self, rel, content):
        p = self.root / rel
        p.parent.mkdir(parents=True, exist_ok=True)
        p.write_text(content, encoding="utf-8")
        return p

    def obs_row(self, ticket):
        return {"session": "ses_" + re.sub(r"\W", "", ticket), "ticket": ticket}

    def test_resolves_by_suffix_under_project_root(self):
        self.write_ticket("ff-wiki-v5/proposals/p/issues/01-x.md", TICKET_S)
        rows = [self.obs_row("proposals/p/issues/01-x.md")]
        strata, unresolved = em.resolve_strata(rows, str(self.root))
        self.assertEqual(strata[rows[0]["session"]], "S")
        self.assertEqual(unresolved, 0)

    def test_missing_ticket_file_is_unknown_and_counted(self):
        rows = [self.obs_row("proposals/p/issues/gone.md"),
                {"session": "ses_noticket", "ticket": "-"}]
        strata, unresolved = em.resolve_strata(rows, str(self.root))
        self.assertEqual(strata["ses_noticket"], "unknown")
        self.assertEqual(unresolved, 2)

    def test_no_project_root_means_all_unknown(self):
        rows = [self.obs_row("proposals/p/issues/01-x.md")]
        strata, unresolved = em.resolve_strata(rows, None)
        self.assertEqual(strata[rows[0]["session"]], "unknown")
        self.assertEqual(unresolved, 1)

    def test_node_modules_pruned_from_glob(self):
        self.write_ticket("node_modules/x/issues/junk.md", TICKET_S)
        self.write_ticket("real/proposals/issues/01-y.md", TICKET_S)
        rows = [self.obs_row("proposals/issues/01-y.md")]
        strata, unresolved = em.resolve_strata(rows, str(self.root))
        self.assertEqual(strata[rows[0]["session"]], "S")
        self.assertEqual(unresolved, 0)


def iso(ms):
    """ISO8601 ts for an epoch-ms fixture value (mirrors extract's fmt_ts)."""
    return datetime.fromtimestamp(ms / 1000, tz=timezone.utc).isoformat()


def rrow(session, ts, dur_min=5.0, model="gemini-3.8-flash-high",
         ticket="proposals/p/issues/01-x.md", steps=10, rep_max=1,
         in_tok=0, out_tok=0, cacheR_tok=0, status="COMPLETED", agent=None):
    """A parsed-observations-row fixture (shape of em.parse_obs_rows output).
    agent=None keeps the key absent — the legacy pre-agent-format row shape;
    passing agent="..." mirrors a row from a multi-agent table."""
    row = {"session": session, "ts": iso(ts), "dur_min": dur_min, "model": model,
           "provider": "prov", "ticket": ticket, "epoch": "-", "steps": steps,
           "tools": steps, "bash_n": 0, "rep_max": rep_max, "fail_streak": 0,
           "in_tok": in_tok, "out_tok": out_tok, "cacheR_tok": cacheR_tok,
           "est_cost": 0.0, "status": status}
    if agent is not None:
        row["agent"] = agent
    return row


class RunawayTests(unittest.TestCase):
    def base(self, **kw):
        d = {"steps": 10, "rep_max": 1, "dur_min": 5.0}
        d.update(kw)
        return d

    def test_at_exact_thresholds_not_runaway(self):
        r = self.base(steps=em.RUNAWAY_STEPS, rep_max=em.RUNAWAY_REP_MAX,
                      dur_min=em.RUNAWAY_DUR_MIN)
        self.assertFalse(em.is_runaway(r))

    def test_each_condition_alone_triggers(self):
        self.assertTrue(em.is_runaway(self.base(steps=em.RUNAWAY_STEPS + 1)))
        self.assertTrue(em.is_runaway(self.base(rep_max=em.RUNAWAY_REP_MAX + 1)))
        self.assertTrue(em.is_runaway(self.base(dur_min=em.RUNAWAY_DUR_MIN + 0.1)))


class ReportCostTests(unittest.TestCase):
    def row(self, model, ti, to, tcr):
        return {"model": model, "in_tok": ti, "out_tok": to, "cacheR_tok": tcr}

    def test_est_cost_recomputed_from_tokens_and_tier(self):
        r = self.row("gemini-3.8-flash-high", 1_000_000, 1_000_000, 1_000_000)
        self.assertEqual(em.est_cost_for(r, em.DEFAULT_PRICES), 2.875)
        r2 = self.row("glm-4.7", 1_000_000, 1_000_000, 1_000_000)
        self.assertEqual(em.est_cost_for(r2, em.DEFAULT_PRICES),
                         round(0.60 + 2.20 + 0.113, 4))

    def test_price_override_replaces_tier_triple(self):
        r = self.row("gemini-3.8-flash-high", 1_000_000, 1_000_000, 1_000_000)
        prices = dict(em.DEFAULT_PRICES)
        prices["flash"] = (1.0, 1.0, 1.0)
        self.assertEqual(em.est_cost_for(r, prices), 3.0)


class AggregateTests(unittest.TestCase):
    EPOCHS = None  # set per test via em.parse_epochs

    def agg(self, rows, epochs_spec="pre:<1000,mid:1000-2000", strata=None):
        epochs = em.parse_epochs(epochs_spec)
        strata = strata or {r["session"]: "L" for r in rows}
        return em.aggregate_cells(rows, epochs, strata, em.DEFAULT_PRICES)

    def test_even_n_median_is_mean_of_two_middle_values(self):
        # out_tok at 2.50/M -> est costs 1.0/2.0/3.0/4.0; median 2.5
        rows = [rrow(f"ses_m{i}", 500, dur_min=float(d),
                     out_tok=int(c * 400_000))
                for i, (d, c) in enumerate([(10.0, 1.0), (20.0, 2.0),
                                            (30.0, 3.0), (40.0, 4.0)])]
        agg = self.agg(rows)
        cell = agg["pre"]["cells"]["L"]
        self.assertEqual(cell["n"], 4)
        self.assertEqual(cell["cost_med"], 2.5)   # (2+3)/2
        self.assertEqual(cell["dur_med"], 25.0)   # (20+30)/2

    def test_cell_runaway_count_and_cache_share(self):
        rows = [
            rrow("ses_a1", 500, out_tok=400_000),                    # cost 1.0
            rrow("ses_a2", 600, cacheR_tok=10_000_000, steps=999),   # cost 0.75, runaway
        ]
        cell = self.agg(rows)["pre"]["cells"]["L"]
        self.assertEqual(cell["runaway"], 1)
        self.assertAlmostEqual(cell["cache_share"], 0.75 / 1.75)

    def test_tier_split_medians_within_cell(self):
        rows = [
            rrow("ses_f1", 500, model="gemini-flash", out_tok=400_000),   # flash 1.0
            rrow("ses_f2", 600, model="gemini-flash", out_tok=800_000),   # flash 2.0
            rrow("ses_g1", 700, model="glm-4.7", out_tok=1_000_000),      # flagship 2.2
        ]
        cell = self.agg(rows)["pre"]["cells"]["L"]
        self.assertEqual(cell["flash_n"], 2)
        self.assertEqual(cell["flash_med"], 1.5)
        self.assertEqual(cell["flagship_n"], 1)
        self.assertEqual(cell["flagship_med"], 2.2)

    def test_status_counts_and_epoch_totals(self):
        rows = [
            rrow("ses_s1", 500, status="COMPLETED"),
            rrow("ses_s2", 500, status="BLOCKED"),
            rrow("ses_s3", 500, status="none"),
            rrow("ses_s4", 1500, status="INCONCLUSIVE"),
        ]
        agg = self.agg(rows)
        self.assertEqual(agg["pre"]["status"],
                         {"COMPLETED": 1, "BLOCKED": 1, "none": 1})
        self.assertEqual(agg["mid"]["status"],
                         {"INCONCLUSIVE": 1})

    def test_unclassified_epoch_bucket(self):
        rows = [rrow("ses_u1", 9999)]
        agg = self.agg(rows, epochs_spec="pre:<1000")
        self.assertIn(em.UNCLASSIFIED, agg)
        self.assertEqual(agg[em.UNCLASSIFIED]["cells"]["L"]["n"], 1)


class RenderTests(unittest.TestCase):
    def render(self, rows, epochs_spec, strata, warnings=(), **kw):
        epochs = em.parse_epochs(epochs_spec)
        prices = kw.pop("prices", em.DEFAULT_PRICES)
        return em.render_report(
            rows, epochs, strata, prices,
            obs_name="observations.md", epochs_spec=epochs_spec,
            warnings=list(warnings), **kw)

    def test_all_epoch_rows_present_including_empty(self):
        rows = [rrow("ses_p1", 500)]
        text = self.render(rows, "pre:<1000,mid:1000-2000,post:>2000",
                           {"ses_p1": "L"})
        self.assertIn("| pre | L | 1 |", text)
        self.assertIn("| mid | - | 0 | - | - | - | 0 |", text)
        self.assertIn("| post | - | 0 | - | - | - | 0 |", text)

    def test_comparison_row_has_both_medians_and_ratio(self):
        rows = [
            rrow("ses_f1", 500, model="gemini-flash", out_tok=400_000),   # flash 1.0
            rrow("ses_g1", 600, model="glm-4.7", out_tok=1_000_000),      # flagship 2.2
        ]
        text = self.render(rows, "pre:<1000", {"ses_f1": "L", "ses_g1": "L"})
        self.assertIn("flash vs flagship", text)
        self.assertIn("| pre | L | 1 | 1.0000 | 1 | 2.2000 | 0.45 |", text)

    def test_warning_lines_and_status_table(self):
        rows = [rrow("ses_w1", 500, status="none"),
                rrow("ses_w2", 500, status="COMPLETED")]
        text = self.render(rows, "pre:<1000", {"ses_w1": "unknown",
                                               "ses_w2": "S"},
                           warnings=["2 session(s) unresolvable -> unknown"])
        self.assertIn("> warning: 2 session(s) unresolvable", text)
        self.assertIn("| pre | S |", text)
        self.assertIn("| pre | unknown |", text)
        # status table counts none explicitly (requirement: 无终态计数)
        self.assertIn("| pre | 2 | 1 | 0 | 1 | 0 |", text)

    def test_header_carries_epochs_prices_and_generated(self):
        rows = [rrow("ses_h1", 500)]
        text = self.render(rows, "pre:<1000", {"ses_h1": "L"},
                           generated="2026-09-13T21:00:00+08:00")
        self.assertIn("generated: 2026-09-13T21:00:00+08:00", text)
        self.assertIn("pre:<1000", text)
        self.assertIn("flash", text)
        self.assertIn("flagship", text)


REAL_EPOCHS = ("pre-hardening:<1789291500000,"
               "provider-switch:1789291500000-1789298700000,hardened:>1789298700000")


class CliReportTests(unittest.TestCase):
    def setUp(self):
        self._tmp = tempfile.TemporaryDirectory()
        self.dir = pathlib.Path(self._tmp.name)
        self.obs = self.dir / "observations.md"
        self.out = self.dir / "report.md"
        self.root = self.dir / "repo"
        self.addCleanup(self._tmp.cleanup)

    def write_obs(self, rows):
        lines = list(em.HEADER_LINES)
        for r in rows:
            lines.append(em.format_row(r))
        self.obs.write_text("\n".join(lines) + "\n", encoding="utf-8")

    def write_ticket_tree(self):
        p = self.root / "ff-wiki-v5/proposals/p/issues/01-x.md"
        p.parent.mkdir(parents=True, exist_ok=True)
        p.write_text(TICKET_L_BY_COUNT, encoding="utf-8")

    def run_report(self, *extra):
        return run_cli(["report", "--obs", str(self.obs)] + list(extra))

    def read_report(self, path=None):
        return (path or self.out).read_text(encoding="utf-8")

    def test_report_writes_default_out_and_replaces_stale_content(self):
        self.write_obs([rrow("ses_c1", 500)])
        self.write_ticket_tree()
        self.out.write_text("# stale report from yesterday\nmust be replaced\n",
                            encoding="utf-8")
        code, out, err = self.run_report("--epochs", "pre:<1000",
                                         "--project-root", str(self.root))
        self.assertEqual(code, 0, err)
        text = self.read_report()
        self.assertNotIn("stale report", text)
        self.assertIn("# Executor value report", text)
        self.assertIn("| pre | L | 1 |", text)
        self.assertTrue(text.endswith("\n"))
        self.assertIn("1 session(s)", out)

    def test_default_out_is_report_md_next_to_obs(self):
        self.write_obs([rrow("ses_c2", 500)])
        self.write_ticket_tree()
        code, _, err = self.run_report("--epochs", "pre:<1000",
                                       "--project-root", str(self.root))
        self.assertEqual(code, 0, err)
        self.assertTrue(self.out.exists())
        self.assertFalse((self.dir / "observations.md_report").exists())

    def test_report_does_not_touch_observations(self):
        self.write_obs([rrow("ses_c3", 500)])
        self.write_ticket_tree()
        before = self.obs.read_bytes()
        self.run_report("--epochs", "pre:<1000", "--project-root", str(self.root))
        self.assertEqual(self.obs.read_bytes(), before)

    def test_bad_epochs_exits_nonzero_without_producing_report(self):
        self.write_obs([rrow("ses_c4", 500)])
        self.write_ticket_tree()
        for bad in ("no-colon-here", "pre:<12ab"):
            code, _, err = self.run_report("--epochs", bad,
                                           "--project-root", str(self.root))
            self.assertEqual(code, 1)
            self.assertIn("error:", err)
            self.assertNotIn("'\''", err)
            self.assertFalse(self.out.exists(), f"report produced for {bad!r}")
        # a pre-existing report is not truncated by a failed run
        self.out.write_text("previous\n", encoding="utf-8")
        code, _, _ = self.run_report("--epochs", "junk")
        self.assertEqual(code, 1)
        self.assertEqual(self.out.read_text(encoding="utf-8"), "previous\n")

    def test_missing_obs_errors_cleanly(self):
        code, _, err = run_cli(["report", "--obs", str(self.dir / "nope.md")])
        self.assertEqual(code, 1)
        self.assertIn("error:", err)

    def test_invalid_price_override_errors_without_report(self):
        self.write_obs([rrow("ses_c5", 500)])
        code, _, err = self.run_report("--price-override", "{not json")
        self.assertEqual(code, 1)
        self.assertIn("error:", err)
        self.assertFalse(self.out.exists())

    def test_no_matching_tickets_all_unknown_with_top_warning(self):
        self.write_obs([rrow("ses_c6", 500), rrow("ses_c7", 1500)])
        empty_root = self.dir / "empty-repo"
        empty_root.mkdir()
        code, _, err = self.run_report("--epochs", "pre:<1000,mid:1000-2000",
                                       "--project-root", str(empty_root))
        self.assertEqual(code, 0, err)
        text = self.read_report()
        self.assertIn("warning:", text)
        self.assertIn("unknown", text)
        self.assertIn("| pre | unknown | 1 |", text)
        self.assertIn("| mid | unknown | 1 |", text)

    def test_no_project_root_given_all_unknown_with_warning(self):
        self.write_obs([rrow("ses_c8", 500)])
        code, _, err = self.run_report("--epochs", "pre:<1000")
        self.assertEqual(code, 0, err)
        text = self.read_report()
        self.assertIn("warning:", text)
        self.assertIn("| pre | unknown | 1 |", text)

    def test_price_override_flows_into_report(self):
        self.write_obs([rrow("ses_c9", 500, model="gemini-flash",
                             out_tok=1_000_000)])  # default flash price -> 2.5
        self.write_ticket_tree()
        code, _, err = self.run_report("--epochs", "pre:<1000",
                                       "--project-root", str(self.root),
                                       "--price-override", '{"flash": [0, 1, 0]}')
        self.assertEqual(code, 0, err)
        text = self.read_report()
        self.assertIn("| pre | L | 1 | 5 | 1.0000 |", text)
        self.assertIn("overridden", text)

    def test_min_replay_incident_grouping(self):
        self.write_obs([
            rrow("ses_incident", em.ts_ms("2026-09-13T12:35:44+08:00"),
                 dur_min=86.4, steps=895, rep_max=184, out_tok=94_875,
                 cacheR_tok=63_657_803),
            rrow("ses_prov1", em.ts_ms("2026-09-13T17:38:09+08:00"),
                 dur_min=8.5, steps=100, rep_max=10),
            rrow("ses_prov2", em.ts_ms("2026-09-13T18:53:58+08:00"),
                 dur_min=0.1, steps=3, rep_max=0),
        ])
        self.write_ticket_tree()
        code, _, err = self.run_report("--epochs", REAL_EPOCHS,
                                       "--project-root", str(self.root))
        self.assertEqual(code, 0, err)
        text = self.read_report()
        pre_line = next(ln for ln in text.splitlines()
                        if ln.startswith("| pre-hardening | L |"))
        self.assertTrue(pre_line.startswith("| pre-hardening | L | 1 | 86.4 |"))
        self.assertTrue(pre_line.endswith("| 1 |"))  # runaway_n = 1 (incident)
        prov_line = next(ln for ln in text.splitlines()
                         if ln.startswith("| provider-switch | L |"))
        self.assertTrue(prov_line.endswith("| 0 |"))  # no runaway sessions
        self.assertIn("| provider-switch | L | 2 |", prov_line)
        self.assertIn("| hardened | - | 0 | - | - | - | 0 |", text)


HARDENED_TS = 1789300000000  # 2026-09-13T19:46:40+08:00 — inside hardened epoch
PRE_TS = 1789274144000      # 2026-09-13T12:35:44+08:00 — incident session (pre)


class GateG1Tests(unittest.TestCase):
    """G1 loop safety (ticket 03; design d-decision-gates): hardened epoch
    runaway events = 0 -> pass; any -> fail; hardened empty -> insufficient-n."""

    def g1(self, rows):
        hrd = em.hardened_rows(rows, em.parse_epochs(REAL_EPOCHS))
        return em.evaluate_g1(hrd)

    def test_hardened_empty_is_insufficient_n(self):
        rows = [rrow("ses_pre", PRE_TS)]  # baseline-only data never gates
        g = self.g1(rows)
        self.assertEqual(g["verdict"], em.GATE_INSUFFICIENT)
        self.assertIn("0 session", g["detail"])

    def test_hardened_runaway_session_fails_with_ids(self):
        rows = [rrow("ses_ok_h", HARDENED_TS, steps=em.RUNAWAY_STEPS,
                     rep_max=em.RUNAWAY_REP_MAX, dur_min=em.RUNAWAY_DUR_MIN),
                rrow("ses_bad_h", HARDENED_TS + 1, steps=em.RUNAWAY_STEPS + 1)]
        g = self.g1(rows)
        self.assertEqual(g["verdict"], em.GATE_FAIL)
        self.assertIn("ses_bad_h", g["detail"])
        self.assertNotIn("ses_ok_h", g["detail"])

    def test_hardened_clean_sessions_pass(self):
        rows = [rrow("ses_ok_h", HARDENED_TS, steps=em.RUNAWAY_STEPS,
                     rep_max=em.RUNAWAY_REP_MAX, dur_min=em.RUNAWAY_DUR_MIN)]
        g = self.g1(rows)
        self.assertEqual(g["verdict"], em.GATE_PASS)
        self.assertIn("1 session", g["detail"])


class GateG2Tests(unittest.TestCase):
    """G2 economics (design d-decision-gates): same-stratum flash cost median
    < flagship median x 0.7 with n>=5 per side, else insufficient-n."""

    # 1 $/M output for both tiers -> est_cost = out_tok/1e6 (exact, price-independent)
    PRICES = {"flash": (0, 1, 0), "flagship": (0, 1, 0)}

    def build_agg(self, flash_out, flagship_out, ts=HARDENED_TS):
        rows = [rrow(f"ses_f{i}", ts, model="gemini-flash", out_tok=t)
                for i, t in enumerate(flash_out)]
        rows += [rrow(f"ses_g{i}", ts, model="glm-4.7", out_tok=t)
                 for i, t in enumerate(flagship_out)]
        strata = {r["session"]: "L" for r in rows}
        epochs = em.parse_epochs(REAL_EPOCHS)
        return em.aggregate_cells(rows, epochs, strata, self.PRICES)

    def evaluate(self, flash_out, flagship_out, **kw):
        per, overall = em.evaluate_g2(self.build_agg(flash_out, flagship_out, **kw))
        self.assertEqual(len(per), 1)  # only stratum L present
        return per[0], overall

    def test_pass_ratio_below_threshold(self):
        # flash n=6, flagship n=5, median ratio 0.5 -> pass
        entry, overall = self.evaluate([500_000] * 6, [1_000_000] * 5)
        self.assertEqual(entry["verdict"], em.GATE_PASS)
        self.assertEqual(overall["verdict"], em.GATE_PASS)
        self.assertIn("ratio 0.50", entry["detail"])

    def test_insufficient_n_when_either_side_below_five(self):
        entry, overall = self.evaluate([500_000] * 6, [1_000_000] * 4)
        self.assertEqual(entry["verdict"], em.GATE_INSUFFICIENT)
        self.assertEqual(overall["verdict"], em.GATE_INSUFFICIENT)
        self.assertIn("flagship n=4", entry["detail"])
        entry, _ = self.evaluate([500_000] * 4, [1_000_000] * 5)
        self.assertEqual(entry["verdict"], em.GATE_INSUFFICIENT)
        self.assertIn("flash n=4", entry["detail"])

    def test_fail_ratio_above_threshold_extends_observation(self):
        entry, overall = self.evaluate([800_000] * 6, [1_000_000] * 5)  # 0.8
        self.assertEqual(entry["verdict"], em.GATE_FAIL)
        self.assertEqual(overall["verdict"], em.GATE_FAIL)
        self.assertIn("extend observation", entry["detail"])
        self.assertIn("ratio 0.80", entry["detail"])

    def test_fail_ratio_at_exactly_0p7_is_not_pass(self):
        # pinned rule is strict: median < flagship x 0.7 (equality -> fail)
        entry, _ = self.evaluate([700_000] * 6, [1_000_000] * 5)  # 0.7 == 0.7
        self.assertEqual(entry["verdict"], em.GATE_FAIL)

    def test_fail_ratio_at_2x_maps_to_rollback(self):
        entry, _ = self.evaluate([2_500_000] * 6, [1_000_000] * 5)  # 2.5x
        self.assertEqual(entry["verdict"], em.GATE_FAIL)
        self.assertIn("rollback", entry["detail"])

    def test_baseline_epoch_data_never_gates(self):
        # identical shape placed in pre-hardening -> nothing judged
        per, overall = em.evaluate_g2(self.build_agg(
            [500_000] * 6, [1_000_000] * 5, ts=PRE_TS))
        self.assertEqual(per, [])
        self.assertEqual(overall["verdict"], em.GATE_INSUFFICIENT)


class GateG3Tests(unittest.TestCase):
    """G3 duration (design d-decision-gates): same-stratum flash dur median
    <= flagship x 1.5 — recorded only, never triggers rollback by itself."""

    def rows_for(self, flash_dur, flag_dur, n=5, st="S"):
        rows = [rrow(f"ses_f{i}", HARDENED_TS, model="gemini-flash",
                     dur_min=flash_dur) for i in range(n)]
        rows += [rrow(f"ses_g{i}", HARDENED_TS, model="glm-4.7",
                      dur_min=flag_dur) for i in range(n)]
        return rows, {r["session"]: st for r in rows}

    def evaluate(self, flash_dur, flag_dur, n=5):
        rows, strata = self.rows_for(flash_dur, flag_dur, n)
        per, overall = em.evaluate_g3(em.hardened_rows(
            rows, em.parse_epochs(REAL_EPOCHS)), strata)
        return per, overall

    def test_ratio_exactly_1p5_passes(self):
        per, overall = self.evaluate(3.0, 2.0)
        self.assertEqual(per[0]["verdict"], em.GATE_PASS)  # <= inclusive
        self.assertEqual(overall["verdict"], em.GATE_PASS)
        self.assertIn("ratio 1.50", per[0]["detail"])

    def test_ratio_above_1p5_fails_recorded_only(self):
        per, overall = self.evaluate(3.1, 2.0)
        self.assertEqual(per[0]["verdict"], em.GATE_FAIL)
        self.assertEqual(overall["verdict"], em.GATE_FAIL)
        self.assertIn("recorded only", per[0]["detail"])

    def test_insufficient_n(self):
        per, overall = self.evaluate(3.0, 2.0, n=4)
        self.assertEqual(per[0]["verdict"], em.GATE_INSUFFICIENT)
        self.assertEqual(overall["verdict"], em.GATE_INSUFFICIENT)


class GateG4Tests(unittest.TestCase):
    """G4 circuit-break (design d-decision-gates): hardened row
    with rep_max>=20 or steps>=400 (right-inclusive) -> fail + rollback alarm."""

    def g4(self, **kw):
        hrd = em.hardened_rows([rrow("ses_t", HARDENED_TS, **kw)],
                               em.parse_epochs(REAL_EPOCHS))
        return em.evaluate_g4(hrd)

    def test_rep_max_20_triggers(self):
        g = self.g4(rep_max=em.G4_REP_MAX)
        self.assertEqual(g["verdict"], em.GATE_FAIL)
        self.assertIn("ses_t", g["detail"])

    def test_rep_max_19_does_not_trigger(self):
        self.assertEqual(self.g4(rep_max=em.G4_REP_MAX - 1)["verdict"],
                         em.GATE_NOT_TRIGGERED)

    def test_steps_400_triggers(self):
        self.assertEqual(self.g4(steps=em.G4_STEPS)["verdict"], em.GATE_FAIL)

    def test_steps_399_does_not_trigger(self):
        self.assertEqual(self.g4(steps=em.G4_STEPS - 1)["verdict"],
                         em.GATE_NOT_TRIGGERED)

    def test_baseline_epoch_shape_never_triggers(self):
        # pre-hardening incident shape (rep_max=184) must NOT trip the gate:
        # only the hardened epoch is judged
        hrd = em.hardened_rows(
            [rrow("ses_incident", PRE_TS, rep_max=184, steps=895)],
            em.parse_epochs(REAL_EPOCHS))
        self.assertEqual(em.evaluate_g4(hrd)["verdict"], em.GATE_NOT_TRIGGERED)


class GateCompositionTests(unittest.TestCase):
    """evaluate_gates: G1 fail / G4 triggered raise rollback alarms; G2/G3
    verdicts alone never do."""

    def compose(self, rows, strata=None):
        epochs = em.parse_epochs(REAL_EPOCHS)
        strata = strata or {r["session"]: "L" for r in rows}
        agg = em.aggregate_cells(rows, epochs, strata, em.DEFAULT_PRICES)
        return em.evaluate_gates(rows, strata, epochs, agg)

    def test_g1_runaway_raises_alarm(self):
        result = self.compose(
            [rrow("ses_bad_h", HARDENED_TS, steps=895, rep_max=184,
                  dur_min=86.4)])
        alarms = " ".join(result["alarms"])
        self.assertIn("G1", alarms)
        self.assertIn("ses_bad_h", alarms)
        self.assertIn("rollback", alarms)

    def test_g4_triggered_raises_immediate_rollback_alarm(self):
        result = self.compose([rrow("ses_cb", HARDENED_TS, rep_max=20)])
        alarms = " ".join(result["alarms"])
        self.assertIn("G4", alarms)
        self.assertIn("ses_cb", alarms)
        self.assertIn("immediately", alarms)

    def test_g3_fail_alone_produces_no_alarm(self):
        rows = [rrow(f"ses_f{i}", HARDENED_TS, model="gemini-flash",
                     dur_min=3.1) for i in range(5)]
        rows += [rrow(f"ses_g{i}", HARDENED_TS, model="glm-4.7",
                      dur_min=2.0) for i in range(5)]
        strata = {r["session"]: "S" for r in rows}
        result = self.compose(rows, strata)
        self.assertEqual(result["alarms"], [])
        gates = {g["gate"]: g["verdict"] for g in result["gates"]}
        self.assertEqual(gates["G3 duration (overall)"], em.GATE_FAIL)

    def test_all_gates_present_in_result(self):
        result = self.compose([rrow("ses_pre_only", PRE_TS)])
        names = [g["gate"] for g in result["gates"]]
        for prefix in ("G1 ", "G2 economics (overall)", "G3 duration (overall)",
                       "G4 "):
            self.assertTrue(any(n.startswith(prefix) for n in names), names)


class ObservationWindowTests(unittest.TestCase):
    """Pre-registered window 2026-09-14 00:00 -> 09-20 23:59 local;
    remaining days = window last day - report run day (ticket pinned)."""

    def test_day_before_window_shows_seven_remaining(self):
        line, remaining = em.observation_window("2026-09-13T21:00:00+08:00")
        self.assertEqual(remaining, 7)
        self.assertIn("7 day(s) remaining", line)

    def test_each_day_of_the_window(self):
        self.assertEqual(em.observation_window("2026-09-14T00:30:00+08:00")[1], 6)
        self.assertEqual(em.observation_window("2026-09-16T12:00:00+08:00")[1], 4)
        self.assertEqual(em.observation_window("2026-09-20T23:00:00+08:00")[1], 0)

    def test_after_window_end_reports_ended(self):
        line, remaining = em.observation_window("2026-09-22T08:00:00+08:00")
        self.assertEqual(remaining, -2)
        self.assertIn("window ended", line)

    def test_line_carries_window_bounds_and_local_offset(self):
        line, _ = em.observation_window("2026-09-13T21:00:00+08:00")
        self.assertIn("2026-09-14 00:00", line)
        self.assertIn("2026-09-20 23:59", line)
        self.assertRegex(line, r"UTC[+-]\d{2}:\d{2}")


class GatesRenderTests(unittest.TestCase):
    """render_report: window line + gates section at the top, rollback alarm
    lines above the gates table, gates before the aggregation."""

    def render(self, rows, strata=None, generated="2026-09-13T21:00:00+08:00"):
        epochs = em.parse_epochs(REAL_EPOCHS)
        strata = strata or {r["session"]: "L" for r in rows}
        return em.render_report(
            rows, epochs, strata, em.DEFAULT_PRICES,
            obs_name="observations.md", epochs_spec=REAL_EPOCHS,
            warnings=[], generated=generated)

    def test_baseline_only_data_shows_insufficient_gates(self):
        text = self.render([rrow("ses_pre", PRE_TS)])
        self.assertIn("- observation window:", text)
        self.assertIn("7 day(s) remaining", text)
        self.assertIn("| G1 loop safety (hardened) | insufficient-n |", text)
        self.assertIn("| G2 economics (overall) | insufficient-n |", text)
        self.assertIn("| G3 duration (overall) | insufficient-n |", text)
        self.assertIn("| G4 circuit-break (hardened) | not-triggered |", text)
        self.assertNotIn("ALARM", text)

    def test_gates_section_sits_above_aggregation(self):
        text = self.render([rrow("ses_pre", PRE_TS)])
        self.assertLess(text.index("## Decision gates"),
                        text.index("## Epoch × stratum aggregation"))

    def test_g1_fail_puts_alarm_at_top(self):
        text = self.render([rrow("ses_bad_h", HARDENED_TS, steps=895,
                                 rep_max=184, dur_min=86.4)])
        self.assertIn("> ALARM:", text)
        self.assertIn("ses_bad_h", text)
        self.assertLess(text.index("ALARM"), text.index("## Decision gates"))
        self.assertIn("| G1 loop safety (hardened) | fail |", text)

    def test_g4_triggered_puts_alarm_at_top(self):
        text = self.render([rrow("ses_cb", HARDENED_TS, rep_max=20)])
        self.assertIn("G4 circuit-break TRIGGERED", text)
        self.assertIn("| G4 circuit-break (hardened) | fail |", text)

    def test_g2_judged_from_report_aggregation(self):
        # flash n=5 cost 1.1 vs flagship n=5 cost 2.2 (ratio 0.5) -> G2 pass
        rows = [rrow(f"ses_f{i}", HARDENED_TS, model="gemini-flash",
                     out_tok=440_000) for i in range(5)]
        rows += [rrow(f"ses_g{i}", HARDENED_TS, model="glm-4.7",
                      out_tok=1_000_000) for i in range(5)]
        text = self.render(rows)
        self.assertIn("| G2 economics L | pass |", text)
        self.assertIn("| G2 economics (overall) | pass |", text)

    def test_g3_fail_present_without_alarm(self):
        rows = [rrow(f"ses_f{i}", HARDENED_TS, model="gemini-flash",
                     dur_min=3.1) for i in range(5)]
        rows += [rrow(f"ses_g{i}", HARDENED_TS, model="glm-4.7",
                      dur_min=2.0) for i in range(5)]
        strata = {r["session"]: "S" for r in rows}
        text = self.render(rows, strata)
        self.assertIn("| G3 duration S | fail |", text)
        self.assertNotIn("ALARM", text)
        self.assertIn("never triggers rollback alone", text)


class GatesCliTests(unittest.TestCase):
    """End to end: the report command emits the gates section."""

    def setUp(self):
        self._tmp = tempfile.TemporaryDirectory()
        self.dir = pathlib.Path(self._tmp.name)
        self.obs = self.dir / "observations.md"
        self.out = self.dir / "report.md"
        self.addCleanup(self._tmp.cleanup)

    def test_report_file_contains_window_and_gates(self):
        lines = list(em.HEADER_LINES)
        lines.append(em.format_row(rrow("ses_pre", PRE_TS)))
        self.obs.write_text("\n".join(lines) + "\n", encoding="utf-8")
        code, _, err = run_cli(["report", "--obs", str(self.obs), "--out",
                                str(self.out), "--epochs", REAL_EPOCHS])
        self.assertEqual(code, 0, err)
        text = self.out.read_text(encoding="utf-8")
        self.assertIn("- observation window:", text)
        self.assertIn("## Decision gates", text)
        self.assertIn("| G1 loop safety (hardened) | insufficient-n |", text)


LEGACY_COLUMNS = tuple(c for c in em.COLUMNS if c != "agent")


def legacy_obs_text(*rows):
    """An old-format observations file: 17-column table without agent."""
    lines = ["# Executor value observations", "",
             "Legacy pre-agent-format table (compat fixture).", "",
             "| " + " | ".join(LEGACY_COLUMNS) + " |",
             "|" + "---|" * len(LEGACY_COLUMNS)]
    for r in rows:
        lines.append(em.format_row(r, LEGACY_COLUMNS))
    return "\n".join(lines) + "\n"


class MultiAgentCollectTests(unittest.TestCase):
    """collect_metrics: agent IN (?,...) filter + agent field per row
    (ticket 01 Change 1; default keeps the implementer-only semantics)."""

    def add_all(self, conn):
        add_session(conn, "ses_ma1", tc=BASE_TS + 10)
        add_session(conn, "ses_ma2", tc=BASE_TS + 20, agent="flowforge-investigator")
        add_session(conn, "ses_ma3", tc=BASE_TS + 30, agent="explore")
        add_session(conn, "ses_ma4", tc=BASE_TS + 40, agent="flowforge-reviewer")
        for sid in ("ses_ma1", "ses_ma2", "ses_ma3", "ses_ma4"):
            add_part(conn, sid, dispatch_and_report_parts()[1])

    def test_agents_filter_returns_each_agents_rows(self):
        conn = make_conn()
        self.add_all(conn)
        rows, corrupt = em.collect_metrics(conn, "tangram-v2",
                                           agents=("flowforge-investigator",
                                                   "explore"))
        self.assertEqual(corrupt, 0)
        self.assertEqual([r["session"] for r in rows], ["ses_ma2", "ses_ma3"])
        self.assertEqual([r["agent"] for r in rows],
                         ["flowforge-investigator", "explore"])

    def test_default_agents_is_implementer_only(self):
        conn = make_conn()
        self.add_all(conn)
        rows, _ = em.collect_metrics(conn, "tangram-v2")
        self.assertEqual([r["session"] for r in rows], ["ses_ma1"])
        self.assertEqual(rows[0]["agent"], "flowforge-implementer")

    def test_empty_agents_matches_nothing(self):
        conn = make_conn()
        self.add_all(conn)
        rows, _ = em.collect_metrics(conn, "tangram-v2", agents=())
        self.assertEqual(rows, [])


class MultiAgentCliTests(unittest.TestCase):
    """extract --agents: multi-agent extraction writes the agent column and
    stays idempotent; an agent with no sessions is an empty success."""

    def setUp(self):
        self._tmp = tempfile.TemporaryDirectory()
        self.dir = self._tmp.name
        self.db = os.path.join(self.dir, "fixture.db")
        self.out = os.path.join(self.dir, "obs.md")
        self.addCleanup(self._tmp.cleanup)

    def populate(self, conn):
        add_session(conn, "ses_ag1", tc=BASE_TS + 10,
                    agent="flowforge-investigator")
        add_session(conn, "ses_ag2", tc=BASE_TS + 20, agent="explore")
        add_session(conn, "ses_ag3", tc=BASE_TS + 30)  # implementer: excluded
        for sid in ("ses_ag1", "ses_ag2", "ses_ag3"):
            for p in dispatch_and_report_parts():
                add_part(conn, sid, p)

    def extract_args(self):
        return ["extract", "--db", self.db, "--project", "tangram-v2",
                "--agents", "flowforge-investigator,explore", "--out", self.out]

    def test_extract_multiple_agents_writes_agent_column(self):
        write_db_file(self.db, self.populate)
        code, _, err = run_cli(self.extract_args())
        self.assertEqual(code, 0, err)
        with open(self.out, encoding="utf-8") as f:
            lines = f.read().splitlines()
        header = next(ln for ln in lines if ln.startswith("| session "))
        self.assertEqual(
            header,
            "| session | agent | ts | dur_min | model | provider | ticket | "
            "epoch | steps | tools | bash_n | rep_max | fail_streak | in_tok | "
            "out_tok | cacheR_tok | est_cost | status |")
        rows = [ln for ln in lines if ln.startswith("| ses_")]
        self.assertEqual(len(rows), 2)
        self.assertIn("ses_ag1", rows[0])
        self.assertIn("| flowforge-investigator |", rows[0])
        self.assertIn("ses_ag2", rows[1])
        self.assertIn("| explore |", rows[1])
        self.assertNotIn("ses_ag3", "\n".join(rows))

    def test_multi_agent_rerun_idempotent(self):
        write_db_file(self.db, self.populate)
        self.assertEqual(run_cli(self.extract_args())[0], 0)
        with open(self.out, encoding="utf-8") as f:
            first = f.read()
        code, _, err = run_cli(self.extract_args())
        self.assertEqual(code, 0, err)
        with open(self.out, encoding="utf-8") as f:
            second = f.read()
        self.assertEqual(first, second, "rerun must not change observations.md")

    def test_agent_with_no_sessions_is_empty_success(self):
        write_db_file(self.db, self.populate)
        code, out, err = run_cli(["extract", "--db", self.db,
                                  "--project", "tangram-v2",
                                  "--agents", "nosuch-agent", "--out", self.out])
        self.assertEqual(code, 0, err)
        with open(self.out, encoding="utf-8") as f:
            rows = [ln for ln in f.read().splitlines()
                    if ln.startswith("| ses_")]
        self.assertEqual(rows, [])
        self.assertIn("0 session(s) matched", out)


class LegacyObsFormatTests(unittest.TestCase):
    """Old-format observations (17 columns, no agent): report backfills
    flowforge-implementer; extract appends in the file's own column format
    so the running executor-value-measurement window stays valid."""

    def setUp(self):
        self._tmp = tempfile.TemporaryDirectory()
        self.dir = pathlib.Path(self._tmp.name)
        self.obs = self.dir / "observations.md"
        self.out = self.dir / "report.md"
        self.addCleanup(self._tmp.cleanup)

    def test_parse_backfills_default_agent(self):
        rows = em.parse_obs_rows(
            legacy_obs_text(rrow("ses_lg1", 500), rrow("ses_lg2", 600)))
        self.assertEqual(len(rows), 2)
        self.assertEqual([r["agent"] for r in rows],
                         ["flowforge-implementer"] * 2)

    def test_parse_keeps_agent_values_when_column_present(self):
        lines = list(em.HEADER_LINES)
        lines.append(em.format_row(rrow("ses_lg3", 500, agent="explore")))
        lines.append(em.format_row(rrow("ses_lg4", 600,
                                        agent="flowforge-investigator")))
        rows = em.parse_obs_rows("\n".join(lines) + "\n")
        self.assertEqual([r["agent"] for r in rows],
                         ["explore", "flowforge-investigator"])

    def test_report_on_legacy_obs_groups_under_default_agent(self):
        self.obs.write_text(
            legacy_obs_text(rrow("ses_lg5", HARDENED_TS, steps=10)),
            encoding="utf-8")
        code, _, err = run_cli(["report", "--obs", str(self.obs), "--out",
                                str(self.out), "--epochs", REAL_EPOCHS])
        self.assertEqual(code, 0, err)
        text = self.out.read_text(encoding="utf-8")
        self.assertIn("- agents (1): flowforge-implementer", text)
        self.assertIn("| G1 loop safety (hardened) | pass |", text)
        self.assertIn("| hardened | unknown | 1 |", text)  # no --project-root
        self.assertNotIn("###", text)  # single agent: no per-agent subsections

    def test_extract_appends_in_legacy_files_own_columns(self):
        db = self.dir / "fixture.db"

        def populate(conn):
            add_session(conn, "ses_lg6", tc=BASE_TS + 10)

        write_db_file(str(db), populate)
        self.obs.write_text(legacy_obs_text(rrow("ses_lg_old", 500)),
                            encoding="utf-8")
        code, _, err = run_cli(["extract", "--db", str(db), "--project",
                                "tangram-v2", "--out", str(self.obs)])
        self.assertEqual(code, 0, err)
        text = self.obs.read_text(encoding="utf-8")
        header = next(ln for ln in text.splitlines()
                      if ln.startswith("| session "))
        self.assertNotIn("agent", header)  # legacy header preserved as-is
        new_line = next(ln for ln in text.splitlines() if "ses_lg6" in ln)
        self.assertEqual(len(em._split_table_line(new_line)),
                         len(LEGACY_COLUMNS))
        rows = em.parse_obs_rows(text)  # appended rows keep the table parseable
        self.assertEqual(len(rows), 2)
        self.assertTrue(all(r["agent"] == "flowforge-implementer"
                            for r in rows))


class PerAgentGateTests(unittest.TestCase):
    """Multi-agent report: gates judged per agent (G1/G4 per role, thresholds
    unchanged); one role's runaway never fails another role's gate."""

    def render(self, rows, strata):
        return em.render_report(
            rows, em.parse_epochs(REAL_EPOCHS), strata, em.DEFAULT_PRICES,
            obs_name="observations.md", epochs_spec=REAL_EPOCHS,
            warnings=[], generated="2026-09-13T21:00:00+08:00")

    def test_g1_fail_of_one_agent_does_not_fail_the_other(self):
        rows = [
            rrow("ses_impl_ok", HARDENED_TS, steps=10, rep_max=1),
            rrow("ses_inv_bad", HARDENED_TS, steps=999, rep_max=30,
                 dur_min=60.0, agent="flowforge-investigator"),
        ]
        strata = {"ses_impl_ok": "L", "ses_inv_bad": "L"}
        text = self.render(rows, strata)
        self.assertIn("### flowforge-implementer", text)
        self.assertIn("### flowforge-investigator", text)
        g1 = [ln for ln in text.splitlines()
              if ln.startswith("| G1 loop safety")]
        self.assertEqual(len(g1), 2)  # one G1 verdict per agent
        self.assertIn("| pass |", g1[0])  # implementer (first appearance) clean
        self.assertIn("ses_inv_bad", g1[1])
        self.assertIn("| fail |", g1[1])
        self.assertIn("> ALARM: [flowforge-investigator] G1 loop safety FAILED",
                      text)
        self.assertNotIn("[flowforge-implementer] G1 loop safety FAILED", text)
        agg_section = text[text.index("## Epoch × stratum aggregation"):
                           text.index("runaway_n:")]
        hardened = [ln for ln in agg_section.splitlines()
                    if ln.startswith("| hardened | L | 1 |")]
        self.assertEqual(len(hardened), 2)  # one aggregation row per agent


class DecisionTemplateTests(unittest.TestCase):
    """DECISION.md closeout template exists with the pinned placeholders
    (ticket 03 Generated artifacts; design d-decision-gates 收口产物)."""

    PATH = (pathlib.Path(__file__).resolve().parents[1]
            / "docs/proposals/executor-value-measurement/DECISION.md")

    def test_template_exists_with_required_placeholders(self):
        text = self.PATH.read_text(encoding="utf-8")
        for option in ("keep-flash", "rollback-flagship", "extend-observation"):
            self.assertIn(option, text)          # three-way conclusion
        for gate in ("G1", "G2", "G3", "G4"):
            self.assertIn(f"| {gate}", text)     # verdict table rows
        self.assertIn("observations.md", text)   # raw-data reference
        self.assertIn("date:", text)             # closeout date field
        self.assertIn("2026-09-20", text)        # window end anchor
        self.assertIn("d-decision-gates", text)  # rule authority reference
        self.assertIn("<!--", text)              # placeholder style


if __name__ == "__main__":
    unittest.main()
