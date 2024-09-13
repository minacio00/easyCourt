from locust import HttpUser, task, between


class BookingUser(HttpUser):
    wait_time = between(1, 5)  # Wait 1-5 seconds between tasks

    @task
    def create_booking(self):
        # Assuming your API endpoint is /bookings
        payload = {
            "user_id": 2,
            "opponent_name": "joaao",
            "timeslot_id": 8,
            "singles_flag": True,
        }
        headers = {"Content-Type": "application/json"}

        with self.client.get(
            "/courts", json=payload, headers=headers, catch_response=True
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Failed to create booking: {response.text}")
