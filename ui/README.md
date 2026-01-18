# Shikshakul {SKL}

> The Next-Gen School Management & IAM Platform.

**Shikshakul** is a unified multi-module monorepo designed to handle the complex interactions between school administrators, teachers, and guardians.

## Architecture

The project is built using **Nx** and **Angular**, separating concerns into three distinct applications that share a common "Brain" (libraries).

### Applications (Portals)

| Application         | Port   | Description                                                      |
| :------------------ | :----- | :--------------------------------------------------------------- |
| **Admin Portal**    | `4200` | Superuser controls. Manage IAM, fees, and school configurations. |
| **Teacher Portal**  | `4201` | Staff controls. Manage gradebooks, attendance, and schedules.    |
| **Guarding Portal** | `4202` | End-user view. View report cards, pay fees, and track progress.  |

### Shared Libraries

- `libs/administrator-portal`: Includes components related to admin portal
- `libs/auth`: Identity & Access Management (Login, Roles, Guards).
- `libs/ui-kit`: The Shikshakul Design System (Buttons, Layouts, Themes).
- `libs/data-access`: Core API services, Interfaces, and State.

---

## Workflow: Adding a New Feature

Follow this standard process to add a new feature (e.g., `feature-staff`) to the Administrator Portal.

### 1. Create the Feature Library

Create a container library for the feature. This keeps the domain logic isolated.

```bash
npx nx g @nx/angular:lib libs/administrator-portal/feature-NAME --importPath=@shikshakul/admin/feature-NAME
```

### 2. Create the Main Page

Create the main "smart" component (the page) that will act as the route entry point. Note: Use --export so it can be used in the routing module.

```bash
npx nx g @nx/angular:component libs/administrator-portal/feature-NAME/src/lib/PAGE-NAME-page --export
```

### 3. Create Sub-Components (Dumb Components)

Create smaller, reusable UI components inside a components folder within the library.

```bash
npx nx g @nx/angular:component libs/administrator-portal/feature-NAME/src/lib/components/COMPONENT-NAME
```

---

## Getting Started

### Prerequisites

- Node.js (v18+)
- NPM

### Installation

```bash
npm install
```

### Running the Ecosystem

To start all three portals simultaneously:

```bash
npm start
```

To start a specific portal:

```bash
npm run start:admin
# OR
npm run start:teacher
# OR
npm run start:guardian
```
