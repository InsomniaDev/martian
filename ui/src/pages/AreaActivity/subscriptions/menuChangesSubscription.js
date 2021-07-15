import { gql } from "@apollo/client";

export const subscriptionForMenu = gql`
subscription menuChanges {
  menuChange {
    active
    areaName
    devices {
      id
      areaName
      integration
      name
      state
      type
      value
    }
  }
}
`;
