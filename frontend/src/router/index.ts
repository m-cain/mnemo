import { createRootRoute, createRouter } from "@tanstack/react-router";
import { lazy } from "react";
import { isAuthenticated } from "../lib/apiClient";

// Lazy-loaded route components
const Login = lazy(() => import("../features/auth/components/LoginPage"));
const Register = lazy(() => import("../features/auth/components/RegisterPage"));
const Dashboard = lazy(
  () => import("../features/dashboard/components/DashboardPage")
);

// Inventory route components
const ItemList = lazy(
  () => import("../features/inventory/components/ItemListPage")
);
const ItemDetail = lazy(
  () => import("../features/inventory/components/ItemDetailPage")
);
const CreateItemForm = lazy(
  () => import("../features/inventory/components/CreateItemForm")
);

// Define layouts
const AppLayout = lazy(() => import("../components/layout/AppLayout"));
const AuthLayout = lazy(() => import("../components/layout/AuthLayout"));

// Root route
export const rootRoute = createRootRoute();

// Auth routes
export const authLayoutRoute = rootRoute.createRoute({
  id: "auth",
  component: AuthLayout,
});

export const loginRoute = authLayoutRoute.createRoute({
  path: "/login",
  component: Login,
});

export const registerRoute = authLayoutRoute.createRoute({
  path: "/register",
  component: Register,
});

// Protected routes
export const appLayoutRoute = rootRoute.createRoute({
  id: "app",
  component: AppLayout,
  beforeLoad: async () => {
    // Check if user is authenticated before allowing access to app routes
    if (!isAuthenticated()) {
      throw new Response("Unauthorized", { status: 401 });
    }
  },
  // Handle unauthorized access by redirecting to login
  errorComponent: () => {
    window.location.href = "/login";
    return null;
  },
});

// Dashboard route
export const dashboardRoute = appLayoutRoute.createRoute({
  path: "/",
  component: Dashboard,
});

// Inventory routes
export const inventoryRoute = appLayoutRoute.createRoute({
  path: "/inventory",
  component: ItemList,
});

export const inventoryDetailRoute = appLayoutRoute.createRoute({
  path: "/inventory/$itemId",
  component: ItemDetail,
});

export const inventoryCreateRoute = appLayoutRoute.createRoute({
  path: "/inventory/new",
  component: CreateItemForm,
});

// Create the route tree using the routes
const routeTree = rootRoute.addChildren([
  authLayoutRoute.addChildren([loginRoute, registerRoute]),
  appLayoutRoute.addChildren([
    dashboardRoute,
    inventoryRoute,
    inventoryDetailRoute,
    inventoryCreateRoute,
  ]),
]);

// Create and export the router
export const router = createRouter({ routeTree });

// Register router types
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
