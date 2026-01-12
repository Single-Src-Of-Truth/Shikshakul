# Shikshakul {SKL} 🎓

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

- `libs/auth`: Identity & Access Management (Login, Roles, Guards).
- `libs/ui-kit`: The Shikshakul Design System (Buttons, Layouts, Themes).
- `libs/data-access`: Core API services, Interfaces, and State.

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
