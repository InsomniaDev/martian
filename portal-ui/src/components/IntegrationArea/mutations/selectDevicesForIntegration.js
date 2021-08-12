import { gql } from "@apollo/client";

export const selectDevicesForIntegration = gql`
  mutation selectDevicesForIntegration($integration: String!, $devices: [String], $addDevices: Boolean!, $automationDevice: Boolean!) {
    selectDevicesForIntegration(
      integration: $integration, 
      devices:$devices, 
      addDevices:$addDevices, 
      automationDevice:$automationDevice
    )
  }
`;
