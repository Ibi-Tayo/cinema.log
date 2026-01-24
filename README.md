# Cinema.log

A full-stack film rating application that helps you organise and understand your personal film preferences through comparative rankings. I've used an ELO-based rating system where you compare films head-to-head which (over time) generates personalised rankings. I've also introduced a film graph in the profile that grows as you review and document films.

## Stack

- **Frontend**: Angular 21 with PrimeNG
- **Backend**: Go 1.24+ with vertical slice architecture
- **Database**: PostgreSQL (Docker Compose)
- **Authentication**: GitHub OAuth with JWT tokens
- **External APIs**: TMDB for film search and metadata

## Key Features

- **Bulk Film Comparisons** — Compare up to 50 films at once with progressive loading
- **ELO Rating System** — Dynamic rankings based on pairwise comparisons
- **Film Search** — Real-time search powered by TMDB API
- **Profile Dashboard** — View recent reviews, top-rated films, and films needing more comparisons
- **Interactive Graph Visualization** — Network graph of film relationships using vis-network
- **Personalised Recommendations** — Multi-step recommendation engine with iterative refinement



---
Maintained by ibitayo
