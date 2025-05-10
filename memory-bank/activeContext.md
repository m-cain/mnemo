# Active Context

This document tracks the current focus of work, recent changes, and immediate next steps. It captures active decisions, considerations, important patterns, and recent learnings or insights.

## Current Focus

- Project planning and architecture design for the Mnemo home inventory tracking system
- Designing a flexible and robust data model to support the project requirements
- Establishing project phases and implementation approach

## Recent Changes

- Created comprehensive project plan with phased implementation approach
- Revised database schema to support arbitrary hierarchies of locations
- Defined component structure for both backend and frontend
- Established timeline and milestones for project development

## Next Steps

- Initialize project repositories and set up development environment
- Begin implementing database schema with migration scripts
- Develop core authentication system
- Create base API structure using Chi router
- Set up React frontend project with Shadcn components

## Active Decisions and Considerations

- Adopted a self-referential location model to allow for arbitrary hierarchical structures
- Chosen direct SQL over ORM for better performance and control
- Decided on a progressive implementation approach, starting with core functionality
- Planning to implement location hierarchy with recursive SQL queries
- Considering materialized paths for optimizing location hierarchy queries

## Important Patterns and Preferences

- Backend using Chi router for RESTful API endpoints
- Service-oriented architecture separating business logic
- Direct use of pgx for database access without ORM layers
- React with TypeScript and Shadcn for frontend development
- Test-driven development approach for core functionalities

## Learnings and Insights

- Self-referential database design is essential for supporting flexible location hierarchies
- Progressive implementation will enable faster iteration and validation
- Clear separation between authentication, inventory, and location services will improve maintainability
- Regular documentation updates in the memory bank will be crucial for project continuity
