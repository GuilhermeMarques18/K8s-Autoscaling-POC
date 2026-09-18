import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:9000';

export const options = {
  vus: 10,
  duration: '2m',
};

export default function () {
  const customerId = Math.floor(Math.random() * 100) + 1;
  const payload = JSON.stringify({
    customerId: customerId,
    productId: Math.floor(Math.random() * 20) + 1,
    quantity: Math.floor(Math.random() * 5) + 1,
  });

  const res = http.post(`${BASE_URL}/orders`, payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  check(res, { 'POST /orders: status é 201': (r) => r.status === 201 });
  sleep(1);
}