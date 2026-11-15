export function validateOrder(order) {
  if (!order || typeof order.sku !== 'string') throw new Error('INVALID_ORDER');
  return order;
}
