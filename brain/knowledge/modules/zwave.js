import { DeviceGraph } from "../redisgraph";

export class ZwaveGraph {
  static async InsertIntoGraph(device) {
    const { nodeId, product, trackStatus, type } = device;
    let devicetype = "";
    if (trackStatus) {
      devicetype = "sensor";
    } else {
      devicetype = "misc";
    }

    await DeviceGraph.query(
      `MERGE (:zwave{nodeId:${nodeId}, name:'${type}', product:'${product}',zwavetype:'${devicetype}'})`
    );
    await DeviceGraph.query(`MERGE (:area{name:'zwave'})`);
    await DeviceGraph.query(`MERGE (:devicetype{name:'${devicetype}'})`);
    await DeviceGraph.query(`
        MATCH (a:zwave{nodeId:${nodeId}}), (b:area{name:'zwave'}), (c:devicetype{name:'${devicetype}'})
        MERGE (b) <-[:ZWAVE_NETWORK]- (a) -[:OF_TYPE]-> (c)
      `);
  }

  static async RetrieveZwaveDataFromGraph() {
    const response = await DeviceGraph.query(
      "MATCH (areaName:area) <-[:ZWAVE_NETWORK]- (device:zwave) -[:OF_TYPE]-> (type) RETURN areaName,device,type"
    );

    const records = [];
    while (response.hasNext()) {
      let record = response.next();
      let areaName = record.get("areaName").properties;
      let device = record.get("device").properties;
      let type = record.get("type").properties;
      device.ready = device.ready === "true";
      device.trackStatus = device.trackStatus === "true";
      device.classes = JSON.parse(device.classes);
      device.areaName = areaName.name;
      device.zwaveType = type.name;
      records.push(device);
    }

    return records;
  }
}
