import { DeviceGraph } from "../redisgraph";

export class LutronGraph {
  static async InsertIntoGraph(devices) {
    for (let device of devices) {
      let {
        Name,
        ID,
        Type,
        Area: { Name: AreaName },
      } = device;
      await DeviceGraph.query(
        `CREATE (:lutron{name:'${Name}',deviceId:${ID}})`
      );
      await DeviceGraph.query(`MERGE (:area{name:'${AreaName}'})`);
      await DeviceGraph.query(`MERGE (:devicetype{name:'${Type}'})`);
      await DeviceGraph.query(`
        MATCH (a:lutron{deviceId:${ID}}), (b:area{name:'${AreaName}'}), (c:devicetype{name:'${Type}'})
        CREATE (b) <-[:RESIDES_IN]- (a) -[:OF_TYPE]-> (c)
      `);
    }
  }

  static async RetrieveLutronDataFromGraph() {
    const response = await DeviceGraph.query(
      "MATCH (areaName:area) <-[:RESIDES_IN]- (device:lutron) -[:OF_TYPE]-> (type:devicetype) RETURN areaName,device,type"
    );

    const records = [];
    while (response.hasNext()) {
      let record = response.next();
      let areaName = record.get("areaName").properties;
      let device = record.get("device").properties;
      let type = record.get("type").properties;
      device.areaName = areaName.name;
      device.type = type.name;
      device.id = device.deviceId;
      records.push(device);
    }

    return records;
  }
}
