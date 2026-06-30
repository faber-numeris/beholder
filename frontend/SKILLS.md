# UI — Conventions & Guidelines

## Tech Stack

- **React** (^19) — UI library
- **TypeScript** (^6) — type-safe development
- **React Hook Form** — form state management and validation
- **React Router** (v7) — client-side routing
- **PrimeReact** (^10) — component library
- **PrimeFlex** — layout and utility classes (no custom CSS)
- **PrimeIcons** — icon set
- **Vite** — bundler and dev server

## Folder Structure

```
src/
├── pages/          # Page-level components (login, dashboard, etc.)
├── components/     # Shared/reusable components
├── hooks/          # Custom React hooks
├── services/       # API clients, auth logic
├── App.tsx         # Root component, routing
└── main.tsx        # Entry point, CSS imports
```

## CSS & Styling

- **Do not write custom CSS files.** Use PrimeFlex utility classes for all layout and spacing.
- Import the PrimeReact theme (`lara-light-amber`) in `main.tsx` — no other global stylesheets.
- Use PrimeReact component `className` props with PrimeFlex classes only.
- Avoid inline styles. Use PrimeFlex utility classes instead.

## PrimeFlex Conventions

- Use `flex`, `flex-column`, `flex-row` for layout.
- Use `align-items-*`, `justify-content-*` for alignment.
- Use `gap-*` for spacing between flex children.
- Use `p-*`, `m-*`, `px-*`, `py-*` for padding and margin.
- Use `w-full`, `max-w-*`, `h-full`, `min-h-screen` for sizing.
- Use `text-*`, `font-*` for typography.
- Use `surface-*` and `text-*` for colors.
- Use the `grid` / `col-*` system for responsive layouts.

## PrimeReact Components

- Use PrimeReact components (`Button`, `Card`, `Avatar`, `Chip`, etc.) instead of native HTML elements when a component exists.
- **Do not use `InputText`, `Password`, or `InputNumber` directly in pages.** Use the custom `Input` component (`component/input/Input.tsx`) which wraps them with the standard float-label, error, and validation pattern.
- Configure component props using PrimeReact API, not CSS overrides.

## React Hook Form

- Define form types explicitly with TypeScript interfaces.
- Use `useForm` with `control`, `handleSubmit`, and `formState`.
- Use `Controller` from react-hook-form, bound through the custom `Input` component. Pass `control`, `name`, `label`, `type`, `rules`, and `errors` as props.
- Enable `mode: 'onChange'` for real-time validation.
- Errors are handled internally by the `Input` component via `fieldState.invalid` and the `p-error` class.

## React Router

- Use `createBrowserRouter` with layout routes for auth guards.
- Use `ProtectedRoute` / `PublicRoute` wrapper components to gate access.
- Use `useNavigate` for imperative navigation.
- Route definitions live in `src/router.tsx`.

## Naming

- Files: `PascalCase` for components, `camelCase` for hooks and utilities.
- Exports: default export for page components, named exports for everything else.

## State Management

- Use React `useState` / `useCallback` for local state.
- Use React context for shared state when necessary.
- No external state management library.

## Build & Test Discipline

1. **Always** run `tsc -b` (or `npm run build`) after any code change to verify the project compiles.
2. **Always** run `npm run lint` before committing.
3. If type-checking fails after a change, assume the new code is the cause and fix it before proceeding.
