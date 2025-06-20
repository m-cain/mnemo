# Active Context

## Current Focus

We are developing an inventory management system with a React frontend and a Go backend. The system allows users to track inventory items, categorize them by types, and assign locations. The frontend uses React with TypeScript, TanStack Query and Router, and Tailwind CSS for styling.

## Recent Changes

- Implemented inventory visualizations with Recharts:

  - Created four chart types for inventory analytics:
    - Item type distribution (pie chart)
    - Location distribution (pie chart)
    - Top 10 items by quantity (bar chart)
    - Low stock items visualization (bar chart)
  - Designed responsive chart components with proper loading and empty states
  - Enhanced the Dashboard page with visualization components
  - Added home selection filtering for the charts
  - Improved the overall dashboard organization and layout
  - Created reusable chart components that can be used across the application
  - Added recent items list to the dashboard

- Extended form validation to authentication forms:

  - Created Zod validation schemas for login and registration forms
  - Integrated react-hook-form with Zod validation for auth forms
  - Redesigned the registration form with a modern card-based layout
  - Added password visibility toggle with eye icons for better user experience
  - Implemented password matching validation for registration
  - Enhanced error handling and feedback for auth forms
  - Used the 21st-dev Magic MCP server to generate a sophisticated registration form UI

- Improved LoginPage UI:

  - Modified global styles in `index.css` to ensure the background fills the entire page and the auth card remains centered.
  - Removed conflicting hardcoded button background styles in `index.css` to fix the invisible button text in light mode.

- Implemented dark mode support and fixed routing import errors:

  - Added `ThemeProvider` component to `main.tsx` to manage theme state and apply the `dark` class to the HTML element.
  - Updated `AppLayout.tsx` to use theme-aware background classes (`bg-background`, `bg-card`).
  - Updated `AuthLayout.tsx` to use theme-aware background and text color classes (`bg-background`, `bg-card`, `text-foreground`, `text-muted-foreground`).
  - Updated `LoginPage.tsx` and `RegisterPage.tsx` to use theme-aware text and background classes (`text-foreground`, `bg-background`).
  - Removed a conflicting `@media (prefers-color-scheme: light)` block from `index.css` that was overriding theme variables.
  - Corrected import paths for `LocationsPage` and `ApiKeysPage` in `router/index.ts` and added placeholder components to resolve TypeScript errors.
  - Dark mode is now working correctly on the authentication pages.

- Implemented form validation with Zod and react-hook-form:

  - Created validation schemas for item forms and quantity adjustments
  - Integrated react-hook-form for better form state management
  - Replaced manual validation with schema-based validation
  - Improved type safety with proper TypeScript integration
  - Enhanced user feedback with field-level error messages
  - Created a new QuantityAdjustmentDialog.zod.tsx component using the new validation pattern
  - Integrated shadcn/ui form components for consistent UI and better validation
  - Added proper FormField components with modern rendering pattern
  - Utilized the 21st-dev Magic MCP server for advanced UI component generation

- Added sorting capabilities to the inventory table on ItemListPage:
  - Implemented column header sorting with visual indicators
  - Support for sorting by name, item type, location, and quantity
  - Toggle between ascending and descending order
  - Visual feedback showing current sort column and direction
- Enhanced the ItemDetailPage with:
  - Card-based layout with responsive grid design
  - Improved visual organization of item information
  - Added quantity adjustment dialog for quick inventory updates
  - Better error state handling and visual feedback
  - Low stock indicator with badge
- Improved the CreateItemForm with:
  - Organized sections for better information hierarchy
  - Enhanced validation with field-level error messages
  - Added a helpful tips card with guidance for users
  - Better form submission feedback and loading states
- Created missing UI components (Input, Textarea) for consistent styling
- Added statistics cards to the ItemListPage to show total items, item types, locations count, and low stock items
- Implemented search functionality with item type filtering
- Enhanced the inventory table with better organization and styling
- **Migrated database migration tool from `go-migrate` to `goose`:** Installed `goose` CLI, created `backend/dbconf.yml`, converted existing migration files to `goose` format, and updated the `Makefile` and documentation.

## UI Component Structure

The inventory management UI now includes:

1. Dashboard:

   - Interactive charts for inventory analytics
   - Home selection controls with visual selection state
   - Four chart types for inventory insights:
     - Item Type Distribution (pie chart)
     - Location Distribution (pie chart)
     - Top 10 Items by Quantity (bar chart)
     - Low Stock Items (bar chart)
   - Recent items list showing the latest inventory additions
   - Quick actions menu for common tasks
   - Responsive layout adapting to different screen sizes
   - Loading and empty states for better user experience

2. ItemListPage:

   - Statistics dashboard showing key metrics (total quantities, item types, locations, low stock)
   - Advanced search bar with clear button functionality
   - Item type filtering using a popover component
   - Responsive inventory table with improved styling
   - Sortable columns with ascending/descending indicators (name, type, location, quantity)
   - Low stock indicators with badges

3. ItemDetailPage:

   - Two-column card layout (3-column on larger screens)
   - Item details card with organized information sections
   - Quantity and location card with adjustment controls
   - Confirm dialog for quantity adjustments with improved validation
   - Consistent error and empty states

4. CreateItemForm:
   - Two-column layout with form and tips cards
   - Organized sections for different types of information
   - Enhanced validation with field-level feedback
   - Clear required field indicators
   - Guidance tips for better user experience
5. Chart Components:
   - `InventoryCharts`: Container component orchestrating all chart components
   - `ItemTypeDistributionChart`: Pie chart showing distribution of item types
   - `LocationDistributionChart`: Pie chart showing distribution by location
   - `QuantityTrendsChart`: Bar chart showing top 10 items by quantity
   - `LowStockChart`: Bar chart for items with low stock levels
   - All charts include proper loading states, empty states, and responsive design

## Next Steps

1. Implement bulk operations (import/export, multi-select delete/edit)
2. Improve the home selection UI (currently using a placeholder)
3. Add pagination for inventory tables with large datasets
4. Enhance analytics with time-based trends and usage patterns
5. Add more comprehensive input validation
6. Implement advanced filtering options (date ranges, multiple criteria)

## Technical Decisions

- Using Zod for form validation to ensure type safety and consistent validation
- Using react-hook-form for form state management and validation integration
- Using shadcn/ui for component library to ensure consistent styling and functionality
- Using relative imports (../../../) instead of alias imports (@/) to avoid path resolution issues. This was addressed by creating and running a script to update imports in the `frontend/src/components/ui` directory.
- Implementing type-safe React event handlers
- Utilized the 21st-dev Magic MCP server for generating sophisticated UI components
- Using FormField rendering pattern with render props for better component isolation
- Updated development process to run backend and frontend independently, using Docker Compose only for the database, and added a Makefile for convenience.
- Migrated database migration tool from `go-migrate` to `goose`.

## Recent Issues and Troubleshooting

- Encountered issues with frontend development server startup related to PostCSS and Tailwind CSS configuration.
- Attempted to fix by installing `tw-animate-css` and `@tailwindcss/postcss`.
- Used Perplexity to research correct Tailwind CSS configuration with Vite, which indicated using `@tailwindcss/vite` and updating `vite.config.ts`.
- Attempted to install `@tailwindcss/vite` and update `vite.config.ts`, but still encountered PostCSS errors related to `@tailwindcss/postcss` being required in `postcss.config.js`.
- **Resolved Tailwind CSS with Vite Configuration:** The user independently resolved the issue by installing `tailwindcss` and `@tailwindcss/vite` and removing the `postcss.config.js` file. The configuration is now handled directly in `vite.config.ts`.
- Encountered TypeScript errors related to TanStack Router's route definition syntax ("Property 'createRoute' does not exist").
- Initially attempted to use `new Route` based on outdated documentation.
- Used Perplexity to find the latest documentation and examples for code-based routing with nested routes.
- **Resolved TanStack Router Syntax:** Updated the routing configuration in `frontend/src/router/index.ts` to use the `createRoute` function for all routes and define nesting using the `getParentRoute` option, based on the latest documentation.

## User Experience Focus

- Making inventory status visible at a glance with the stats cards
- Simplifying search, filtering, and sorting operations
- Ensuring responsive design works across device sizes
- Using clear visual indicators for item status (low stock, sorting direction)
- Providing contextual error messages with field-level validation
- Using consistent UI patterns for forms and interactive elements
