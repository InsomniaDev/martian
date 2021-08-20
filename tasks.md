# current work
- hass
  - Create edit screen, and allow for updating variables


## Need to do
- [ ] Need to add unit test coverage
- [ ] Integrate with BoltDB for the backend rather than using the YAML, then setup through APIs implementing different integrations
- [ ] Setup a user uuid for who created the message, we can search by user uuid in the future as well


## Setup Web API
- [ ] Add base web api for consuming incoming calls
- [ ] Add base authentication
- [ ] Add caching layer with a UUID that we can use to tie calls and communication together
- [ ] Need to decide what incoming calls are going to look like

## Brain and Graph data
- Keys between two instances will be concatenated together with a determined symbol
  - For example: `light.office_main&&&light.office_sconces`
  - They will need to be sorted alphabetically so that the alphabetical one appears first
  - The value for this will be the weight that the two of them have
  - This will provide in the future for there to also be more stored in the value
    - Possible things would be what state each device is in and the weight for those 
      - ie. When the `office_main` is `on`, then the `office_sconces` are `off`