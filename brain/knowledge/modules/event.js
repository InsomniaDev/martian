import { EventGraph } from "../redisgraph";

// Get back the events where the weight property is greatest
// match (a:event) -[b:RELATES]-> (c:event) where b.weight > 1 return a,c

export class EventsGraph {
  static async InsertIntoGraph(first) {
    await EventGraph.query(first);
  }

  static async GetRelatedEvents() {
    const response = await EventGraph.query(
      `MATCH (a:event) -[b:RELATES]-> (c:event) WHERE b.weight >= 3 AND b.weight <= 20 AND NOT c.device = "zwave" AND NOT a.deviceId = c.deviceId return a,b,c`
    );

    const records = [];
    while (response.hasNext()) {
      let record = response.next();
      let triggerDevice = record.get("a").properties;
      let weight = record.get("b").properties;
      let nextDevice = record.get("c").properties;
      let recordData = triggerDevice;
      recordData.nextDevice = nextDevice;
      recordData.weight = weight;
      records.push(recordData);
    }
    
    const cleanData = [];
    for (let x of records) {
      if (x.state !== "NaN" && x.nextDevice.state !== "NaN") {
        cleanData.push(x);
      }
    }

    const duplicates = [];
    for (let i = 0; i < cleanData.length; i++) {
      const newElement = cleanData[i];
      cleanData.splice(i, 1);
      newElement.nextDevice.nextDevice = findNestedElements(
        newElement,
        cleanData
      );
      duplicates.push(newElement);
    }
    console.log(JSON.stringify(duplicates));
  }
}

function findNestedElements(device, data) {
  const upperWeight = device.weight.weight + 1;
  const lowerWeight = device.weight.weight - 1;
  for (let x = 0; x < data.length; x++) {
    // Check if the weighting is close enough to be similar
    if (
      device.nextDevice.deviceId === data[x].deviceId &&
      device.nextDevice.device === data[x].device &&
      data[x].weight.weight >= lowerWeight &&
      data[x].weight.weight <= upperWeight
    ) {

      let newData = JSON.parse(JSON.stringify(data));
      newData.splice(x, 1);
      return {
        device: data[x],
        nextDevice: findNestedElements(data[x], newData),
      };
    }
  }
  return null;
}
