import { gql } from "@apollo/client";

export const changeDeviceStatus = gql`
  mutation changeDeviceStatus (
    $id: String!, 
    $integration: String!, 
    $status: String!, 
    $level: String!) {
    changeDeviceStatus(
      id: $id, 
      integration: $integration, 
      status: $status, 
      level: $level)
  }
`;