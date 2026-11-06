import { validateOrder } from './contract.mjs';
const prices = { A: 125, B: 300 };
export function quote(order) {
  validateOrder(order);
  if (!(order.sku in prices)) throw new Error('UNKNOWN_SKU');
  return prices[order.sku];
}
