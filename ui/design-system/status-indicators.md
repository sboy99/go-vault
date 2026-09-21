# Status indicators

| Status | Color token | UI |
|---|---|---|
| success / succeeded / ok | `--ui-success` | Soft badge + check icon |
| failed / error | `--ui-error` | Soft badge + X icon |
| running / queued / info | `--ui-info` | Soft badge + activity icon |
| warning | `--ui-warning` | Soft badge |

Do not repaint success/error to the primary accent — semantic colors stay fixed so operators can scan status quickly while the accent remains for actions and navigation.
