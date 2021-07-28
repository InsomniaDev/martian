import { gql } from "@apollo/client";

// const rokuActivity = "40771265";
export const getIntegrations = gql`
query integrations {
    integrations {
      integrations
      harmony {
        activityId
        id
        name
      }
      hass {
        devices {
          areaName
          name
          entityId
          state
          type
          group
        }
        selectedDevices {
          areaName
          name
          entityId
          state
          type
          group
        }
        token
        url
      }
      kasa {
        plugs {
          areaName
          id
          ipAddress
          name
          type
        }
      }
      lutron {
        config {
          port
          url
          password
          username
          file
        }
        inventory {
          areaName
          id
          name
          state
          type
          value
        }
      }
    }
  }
`;