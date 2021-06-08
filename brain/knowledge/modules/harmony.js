import { DeviceGraph } from "../redisgraph";

export class HarmonyGraph {
  static async InsertIntoGraph(devices) {
    for (let device of devices) {
      let { activityId, name, actions } = device;
      await DeviceGraph.query(
        `CREATE (:activity{name:'${name}',deviceId:'${activityId}', activityId: '${activityId}', actions: '${JSON.stringify(
          actions
        )}'})`
      );
      await DeviceGraph.query(`MERGE (:devicetype{name:'harmony'})`);
      await DeviceGraph.query(`
        MATCH (a:activity{deviceId:'${activityId}'}), (c:devicetype{name:'harmony'})
        CREATE (a) -[:OF_TYPE]-> (c)
      `);
    }
  }

  static async RetrieveLutronDataFromGraph() {
    const response = await DeviceGraph.query(
      "MATCH (device:activity) -[:OF_TYPE]-> (:devicetype{name:'harmony'}) RETURN device"
    );

    const records = [];
    while (response.hasNext()) {
      let record = response.next();
      let devices = record.get("device").properties;
      devices.id = devices.deviceId;
      try {
        if (
          devices.actions &&
          devices.actions !== undefined &&
          devices.actions !== "undefined"
        ) {
          devices.actions = JSON.parse(devices.actions);
        }
      } catch (e) {
        console.log(devices);
      }
      records.push(devices);
    }

    return records;
  }
}
