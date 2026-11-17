import http from 'k6/http';
import { check } from 'k6';

export const options = {
    scenarios: {
        race_condition_test: {
            executor: 'per-vu-iterations',
            vus: 100,
            iterations: 1,
            maxDuration: '10s'
        },
    },
};

export default function () {
    const url = 'http://app:8080/bookings';

    const payload = JSON.stringify({
        event_id: 1,
        seats: [
            {
                row: 1,
                column: 1
            }
        ]
    });

    // 1. Вставьте выведенный токен из $res.token в эту переменную:
    const token = 'YOUR_JWT_TOKEN';

    const params = {
        headers: {
            'Content-Type': 'application/json',
            // 2. Токен подставляется в заголовок Authorization
            'Authorization': `Bearer ${token}`,
        },
    };

    const res = http.post(url, payload, params);

    if (__VU === 1) {
        console.log(`Status: ${res.status} | Body: ${res.body}`);
    }

    check(res, {
        'Status is 200/201 (Success)': (r) => r.status === 200 || r.status === 201,
        'Status is expected rejection (400/403/409/404)': (r) => [400, 403, 409, 404].includes(r.status),
        'No Server Error (500)': (r) => r.status !== 500,
    });
}