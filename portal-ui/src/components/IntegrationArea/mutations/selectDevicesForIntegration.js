import { gql } from "@apollo/client";

export const selectDevicesForIntegration = gql`
  mutation selectDevicesForIntegration($integration: String!, $devices: [String], $addDevices: Boolean!) {
    selectDevicesForIntegration(integration: $integration, devices:$devices, addDevices:$addDevices)
  }
`;
