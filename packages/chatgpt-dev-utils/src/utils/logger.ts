import signale from "signale";


export const logger = new signale.Signale({
  scope: "cgpt",
  config: {
    displayTimestamp: true,
    displayScope: true,
    displayBadge: true,
    displayLabel: true
  }
});
