import {
  createRootRoute,
  createRoute,
  createRouter,
} from "@tanstack/react-router";
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

// Define routes
const rootRoute = createRootRoute();

const authLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: "auth",
  component: AuthLayout,
});

const loginRoute = createRoute({
  getParentRoute: () => authLayoutRoute,
  path: "/login",
  component: Login,
});

const registerRoute = createRoute({
  getParentRoute: () => authLayoutRoute,
  path: "/register",
  component: Register,
});

const appLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
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

const dashboardRoute = createRoute({
  getParentRoute: () => appLayoutRoute,
  path: "/",
  component: Dashboard,
});

const inventoryRoute = createRoute({
  getParentRoute: () => appLayoutRoute,
  path: "/inventory",
  component: ItemList,
});

const inventoryDetailRoute = createRoute({
  getParentRoute: () => appLayoutRoute,
  path: "/inventory/$itemId",
  component: ItemDetail,
});

const inventoryCreateRoute = createRoute({
  getParentRoute: () => appLayoutRoute,
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
