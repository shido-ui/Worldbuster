# World Map & Travel

Phase 5 establishes the mobile-friendly world navigation foundation.

The map is data-driven: locations and routes are domain data rather than UI hard-coding. Travel is server-authoritative and time-based.

Initial districts:
- Central District
- Harbor Ward
- Old Town
- Industrial Belt
- North Highlands

Rules:
- A character can have one active travel operation.
- Travel duration comes from the authoritative route table.
- Clients request travel; the server validates the route and creates the arrival time.
- A client must not decide when arrival occurs.
- Completed travel is resolved from server time.
- Future map expansion can add regions/routes without rewriting the Android navigation layer.

The same route data can feed Android map UI, NPC movement, missions, businesses and later dynamic world events.
