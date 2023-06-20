import http from "k6/http";
import { sleep } from "k6";

export const options = {
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<100'],
  },
};

export default function () {
  http.get("http://test.k6.io");
  sleep(1);
}
