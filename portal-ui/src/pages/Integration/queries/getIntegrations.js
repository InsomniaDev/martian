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
      ipAddress
      areaName
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
      interfaceDevices
      automatedDevices
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
      devices {
        areaName
        id
        name
        state
        type
        value
        lutronName
      }
      interfaceDevices 
      automatedDevices 
    }
  }
}
`;