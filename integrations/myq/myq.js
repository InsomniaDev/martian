import axios from "axios";
import { Config } from "../../src/models/config";

// TODO: Sync with the graph database so that we don't need to keep making calls for the garage door

const constants = {
  endpoint: "https://myqexternal.myqdevice.com",
  appId: "NWknvuBd7LoFHfXmKNMBcgajXtZEgKUh4V7WNzMidrpUUluDpVYVZx+xT4PCM5Kx",
  allTypeIds: [1, 2, 3, 5, 7, 9, 13, 15, 16, 17],
  errorMessages: {
    11: "Something unexpected happened. Please wait a bit and try again.",
    12: "MyQ service is currently down. Please wait a bit and try again.",
    13: "Not logged in.",
    14: "Email and/or password are incorrect.",
    15: "Invalid parameter(s) provided.",
    16: "User will be locked out due to too many tries. 1 try left.",
    17: "User is locked out due to too many tries. Please reset password and try again.",
  },
  doorStates: {
    1: "open",
    2: "closed",
    3: "stopped in the middle",
    4: "going up",
    5: "going down",
    9: "not closed",
  },
  lightStates: {
    0: "off",
    1: "on",
  },
  types: {
    1: "Gateway",
    2: "GDO",
    3: "Light",
    5: "Gate",
    7: "VGDO_Garage_Door",
    9: "Commercial_Door_Operator",
    13: "Camera",
    15: "WGDO_Gateway_AC",
    16: "WGDO_Gateway_DC",
    17: "WGDO_Garage_Door",
  },
  garageDoor: {
    close: 0,
    open: 1,
  },
};

export class MyqGarage {
  constructor() {
    const { username, password } = new Config().getMyqConfig();
    this.username = username;
    this.password = password;
    this.devices = [];
  }

  async login() {
    if (!this.username || !this.password) {
      return this.returnError(14);
    }

    try {
      let response = await axios({
        method: "post",
        url: `${constants.endpoint}/api/v4/User/Validate`,
        headers: {
          MyQApplicationId: constants.appId,
        },
        data: {
          username: this.username,
          password: this.password,
        },
      });
      if (!response || !response.data) {
        return this.returnError(12);
      }
      const { data } = response;
      if (data.ReturnCode === "203") {
        return this.returnError(14);
      } else if (data.ReturnCode === "205") {
        return this.returnError(16);
      } else if (data.ReturnCode === "207") {
        return this.returnError(17);
      } else if (!data.SecurityToken) {
        return this.returnError(11);
      }
      this.securityToken = data.SecurityToken;
      const result = {
        returnCode: 0,
        token: data.SecurityToken,
      };

      // Get all of the devices for the MyQ integration and store them locally
      await this.getDevices();

      return result;
    } catch (err) {
      if (err.statusCode === 500) {
        return this.returnError(14);
      }
      return this.returnError(11, err);
    }
  }

  async getDevices() {
    let response = await axios({
      method: "get",
      url: `${constants.endpoint}/api/v4/userdevicedetails/get`,
      headers: {
        MyQApplicationId: constants.appId,
        securityToken: this.securityToken,
      },
    });

    if (!response || !response.data) {
      return this.returnError(12);
    }

    const { data } = response;
    const { Devices } = data;

    if (data.ReturnCode === "-3333") {
      return this.returnError(13);
    } else if (!Devices) {
      return this.returnError(11);
    }

    for (let device of Devices) {
      device.name = `MyQ_${device.MyQDeviceId}`;
      device.device_type = constants.types[device.MyQDeviceTypeId];
      this.devices.push(device);
      await this.getDoorState(device.MyQDeviceId);
    }
  }

  async getDoorState(id) {
    let result = await this.getDeviceState(this.securityToken, id, "doorstate");

    for (let device of this.devices) {
      if (device.MyQDeviceId === id) {
        device.state = result.state;
        device.stateDescription = constants.doorStates[result.state];
      }
    }
    return { state: constants.doorStates[result.state] };
  }

  async closeGarageDoor(id) {
    await this.setDoorState(id, constants.garageDoor.close);
    await this.getDoorState(id);
  }

  async openGarageDoor(id) {
    await this.setDoorState(id, constants.garageDoor.open);
    await this.getDoorState(id);
  }

  async getLightState(id) {
    let result = await this.getDeviceState(
      this.securityToken,
      id,
      "lightstate"
    );

    if (result.returnCode !== 0) {
      return result;
    }

    const newResult = JSON.parse(JSON.stringify(result));
    newResult.lightState = newResult.state;
    newResult.lightStateDescription =
      constants.lightStates[newResult.lightState];
    delete newResult.state;
    return newResult;
  }

  async setDoorState(id, toggle) {
    let result = await this.setDeviceState(
      this.securityToken,
      id,
      toggle,
      "desireddoorstate"
    );
    return result;
  }

  async setLightState(id, toggle) {
    let result = await this.setDeviceState(
      this.securityToken,
      id,
      toggle,
      "desiredlightstate"
    );
    return result;
  }

  // Format the error to be printed out
  returnError(returnCode, err) {
    const result = {
      returnCode,
      message: constants.errorMessages[returnCode],
    };
    if (err) {
      result.unhandledError = err;
    }
    console.log(JSON.stringify(result));
    return result;
  }

  async getDeviceState(securityToken, id, attributeName) {
    if (!securityToken) {
      return this.returnError(13);
    }

    const resp = await new Promise((resolve, reject) => {
      axios({
        method: "get",
        url: `${constants.endpoint}/api/v4/deviceattribute/getdeviceattribute`,
        headers: {
          MyQApplicationId: constants.appId,
          SecurityToken: securityToken,
        },
        params: {
          MyQDeviceId: id,
          AttributeName: attributeName,
        },
      })
        .then((response) => {
          if (!response || !response.data) {
            return reject(this.returnError(12));
          }

          const { data } = response;

          if (data.ReturnCode === "-3333") {
            return reject(this.returnError(13));
          } else if (!data.ReturnCode) {
            return reject(this.returnError(11));
          } else if (!data.AttributeValue) {
            return reject(this.returnError(15));
          }

          const state = parseInt(data.AttributeValue, 10);
          const result = {
            returnCode: 0,
            state,
          };
          return resolve(result);
        })
        .catch((err) => {
          if (err.statusCode === 400) {
            return reject(this.returnError(15));
          }

          return reject(this.returnError(11, err));
        });
    }).catch((err) => {
      console.log(err);
    });
    return resp;
  }

  async setDeviceState(securityToken, id, toggle, attributeName) {
    if (!securityToken) {
      return returnError(13);
    } else if (toggle !== 0 && toggle !== 1) {
      return returnError(15);
    }

    let response = await axios({
      method: "put",
      url: `${constants.endpoint}/api/v4/deviceattribute/putdeviceattribute`,
      headers: {
        MyQApplicationId: constants.appId,
        SecurityToken: securityToken,
      },
      data: {
        MyQDeviceId: id,
        AttributeName: attributeName,
        AttributeValue: toggle,
      },
    });

    if (!response || !response.data) {
      return reject(this.returnError(12));
    }
    const { data } = response;

    if (data.ReturnCode === "-3333") {
      return reject(this.returnError(13));
    } else if (!data.ReturnCode) {
      return reject(this.returnError(11));
    }

    const result = {
      returnCode: 0,
    };
    return result;
  }
}
