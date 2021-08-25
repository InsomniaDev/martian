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
        interfaceDevices {
          areaName
          name
          entityId
          state
          type
          group
        }
        automatedDevices {
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
        devices {
          areaName
          id
          ipAddress
          name
          type
        }
        interfaceDevices 
        automatedDevices 
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