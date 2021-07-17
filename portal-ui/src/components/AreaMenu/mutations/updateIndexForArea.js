import { gql } from "@apollo/client";

export const updateIndexForArea = gql`
  mutation updateIndexForArea($areaName: String!, $index: Int!) {
    updateIndexForArea(areaName: $areaName, index:$index)
  }
`;