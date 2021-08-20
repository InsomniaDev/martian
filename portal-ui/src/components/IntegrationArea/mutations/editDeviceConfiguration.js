import { gql } from "@apollo/client";

export const editDeviceConfiguration = gql`
  mutation editDeviceConfiguration($integration: String!, $device: String!, $removeEdit: Boolean!) {
    editDeviceConfiguration(
      integration: $integration, 
      device: $device,
      removeEdit: $removeEdit
    )
  }
`;
