import { gql } from "@apollo/client";

export const selectDevicesForIntegration = gql`
  mutation selectDevicesForIntegration($integration: String!, $devices: [String]) {
    selectDevicesForIntegration(integration: $integration, devices:$devices)
  }
`;