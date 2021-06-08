package graphql

// Schema is the graphql schema being used
const Schema = `
  type Query {
    lutronDevice(id: ID!): Lutron
    lutronDevices: [Lutron]
    garageDoor(id: String!): MyQ
    garageDoors: [MyQ]
    checkGarageDoor(id: Int!): GarageDoorState
    currentHarmonyActivity: Harmony
    getHarmonyActivities: [Harmony]
    getHarmonyDevices: [Harmony]
    getHarmonyDeviceCommands(id: String!): [Harmony]
    zwaveDevices: [ZWave]
    kasaDevices: [Kasa]
  }

  type Mutation {
    turnDeviceOn(id: ID!): Lutron
    turnDeviceOff(id: ID!): Lutron
    setLutronDeviceToLevel(id: ID!, level: Int!): Lutron
    myqLogin: [MyQ]
    openGarageDoor(id: Int!): MyQ
    closeGarageDoor(id: Int!): MyQ
    startHarmonyActivity(id: String!): Boolean
    sendHarmonyCommand(
      command: String!
      type: String!
      deviceId: String!
    ): Boolean
    turnAllLightsOn: Boolean
    turnAllLightsOff: Boolean
    healZwaveNetwork: Boolean
    turnOnKasaSwitch(id: ID!): Boolean
    turnOffKasaSwitch(id: ID!): Boolean
  }

  type Subscription {
    lutronChange(id: Int!): Lutron
    lutronChanges: [Lutron]
    harmonyChange: Harmony
    anyLightOn: Boolean
  }

  type Kasa {
    name: String
    host: String
    deviceId: String
    areaName: String
  }

  type Harmony {
    activityId: String
    name: String
    actions: [HarmonyAction]
  }

  type HarmonyAction {
    label: String
    action: HarmonyCommand
  }

  type HarmonyCommand {
    command: String
    type: String
    deviceId: String
  }

  type Lutron {
    _id: String!
    name: String!
    id: ID!
    type: String!
    buttons: [String]
    areaName: String!
    state: String
    value: Float
  }

  type ZWave {
    nodeId: String
    manufacturer: String
    manufacturerid: String
    product: String
    producttype: String
    productid: String
    type: String
    name: String
    loc: String
    classes: String
    ready: Boolean
    trackStatus: Boolean
  }

  type GarageDoorState {
    state: String
  }

  type MyQ {
    name: String
    device_type: String
    state: Int
    stateDescription: String
    MyQDeviceId: Int
    ParentMyQDeviceId: Int
    MyQDeviceTypeId: Int
    MyQDeviceTypeName: String
    RegistrationDateTime: String
    SerialNumber: String
    UserName: String
    UserCountryId: Int
    ChildrenMyQDeviceIds: String
    UpdatedBy: String
    UpdatedDate: String
    ConnectServerDeviceId: String
    Attributes: [MyQAttributes]
  }

  type MyQAttributes {
    MyQDeviceTypeAttributeId: Int
    Value: String
    UpdatedTime: String
    IsDeviceProperty: Boolean
    AttributeDisplayName: String
    IsPersistent: Boolean
    IsTimeSeries: Boolean
    IsGlobal: Boolean
    UpdatedDate: String
  }
`
