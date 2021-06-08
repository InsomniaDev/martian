import { DeviceGraph } from "../redisgraph";

export class KasaGraph {
  static async InsertIntoGraph(device) {
    // for (let device of devices) {
    let { name, host, deviceId } = device;
    console.log(`${name} ${host} ${deviceId}`);
    console.log(
      `MERGE (:kasa{name:'${name}',deviceId:'${deviceId}',host:'${host}'})`
    );
    await DeviceGraph.query(
      `MERGE (:kasa{name:'${name}',deviceId:'${deviceId}',host:'${host}'})`
    );
    await DeviceGraph.query(`MERGE (:devicetype{name:'kasa'})`);
    await DeviceGraph.query(`
        MATCH (a:kasa{deviceId:'${deviceId}'}), (c:devicetype{name:'kasa'})
        CREATE (a) -[:OF_TYPE]-> (c)
      `);
    // }
  }

  static async RetrieveZwaveDataFromGraph() {
    const response = await DeviceGraph.query(
      "MATCH (device:kasa) -[:OF_TYPE]-> (type:devicetype) RETURN device"
    );
    const records = [];
    while (response.hasNext()) {
      let record = response.next();
      let device = record.get("device").properties;
      records.push(device);
    }

    return records;
  }

  static async AddArea(deviceId, area) {
    await DeviceGraph.query(
      `MATCH (device:kasa{deviceId:'${deviceId}'})
       MERGE (area:area{name:'${area}'})
       CREATE (area) <-[:RESIDES_IN]- (device)`
    );
  }
}
