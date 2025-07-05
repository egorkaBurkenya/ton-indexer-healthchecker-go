# Product Context: TON Indexer Health-Checker

## 1. The "Why"

We are building this utility to solve a critical reliability problem in our TON blockchain data indexing pipeline. The `index-worker` is the heart of our system, responsible for feeding fresh, on-chain data to our marketplace. When it silently fails, our product stops reflecting the real-time state of the blockchain, which directly impacts user trust and the usability of our platform. Manual restarts are not a scalable or reliable solution.

## 2. Problem It Solves

- **Eliminates Silent Failures:** Addresses the core issue of the `index-worker` "hanging" without crashing, a state invisible to standard process supervisors.
- **Automates Recovery:** Removes the need for human intervention to detect and resolve indexer stalls, reducing downtime from hours to minutes.
- **Increases Data Freshness:** Ensures that the data presented to users on the marketplace is consistently up-to-date.
- **Improves System Stability:** Provides a robust, automated mechanism for maintaining the health of a critical infrastructure component.

## 3. User Experience Goals

The "user" of this utility is the DevOps team and the platform itself. The experience should be:

- **Transparent:** The health checker should log its status clearly to `stdout` (for success) and `stderr` (for failure), making debugging straightforward.
- **Seamless:** It must integrate perfectly with Docker Swarm's healthcheck mechanism without requiring complex setup.
- **Reliable:** The checker itself must be lightweight and extremely stable, as it is a key part of the recovery system.
- **Configurable:** All critical parameters must be configurable via environment variables, which is the standard for containerized applications.
