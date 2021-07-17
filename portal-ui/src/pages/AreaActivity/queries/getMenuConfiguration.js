import { gql } from "@apollo/client";

// const rokuActivity = "40771265";
export const getMenuConfiguration = gql`
  query {
    menuConfiguration {
      active
      areaName
      index
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