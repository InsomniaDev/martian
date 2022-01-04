import { makeStyles } from '@material-ui/core/styles';
import React, { useState } from "react";
import classNames from 'classnames';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import Typography from '@material-ui/core/Typography';
import Chip from '@material-ui/core/Chip';
import Button from '@material-ui/core/Button';
import { useMutation } from '@apollo/client';
import { selectDevicesForIntegration } from '../mutations/selectDevicesForIntegration';
import KasaEditMenu from './KasaEditMenu';


const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
        padding: "1rem",
        margin: "3%",
        textAlign: 'center',
        backgroundColor: "rgb(199, 216, 223)",
        width: '90%',
    },
    heading: {
        fontSize: theme.typography.pxToRem(23),
    },
    secondaryHeading: {
        fontSize: theme.typography.pxToRem(20),
        color: theme.palette.text.secondary,
    },
    columnHeading: {
        fontSize: theme.typography.pxToRem(15),
        color: theme.palette.text.secondary,
    },
    deviceHeading: {
        fontSize: theme.typography.pxToRem(12),
        color: theme.palette.text.secondary,
        display: 'block',
    },
    icon: {
        verticalAlign: 'bottom',
        height: 20,
        width: 20,
    },
    details: {
        alignItems: 'center',
    },
    column: {
        flexBasis: '33.33%',
    },
    helper: {
        borderLeft: `2px solid ${theme.palette.divider}`,
        padding: `${theme.spacing()}px ${theme.spacing(2)}px`,
    },
    link: {
        color: theme.palette.primary.main,
        textDecoration: 'none',
        '&:hover': {
            textDecoration: 'underline',
        },
    },
    deviceDetails: {
        display: 'block',
        marginTop: theme.spacing(2),
    },
    buttonDiv: {
        display: 'flex',
        flexDirection: 'row',
    },
    button: {
        display: 'block',
        margin: theme.spacing(1),
        backgroundColor: "#9cdbee",
    },
    formControl: {
        margin: theme.spacing(1),
        minWidth: 400,
    },
    em: {
        fontWeight: "bold",
        fontStyle: "normal"
    },
}));



export function KasaIntegration({ integration, refetchData, areaData }) {
    const classes = useStyles();

    const [selectedDevice, changeSelectedDevice] = useState("");
    const [selectDevicesForIntegrationMutation] = useMutation(selectDevicesForIntegration);
    const handleChange = (event) => {
        changeSelectedDevice(event.target.value);
    };

    // Add the selected variable to the interfaceDevices for kasa
    const addToSelectedInterface = () => {
        devices = [selectedDevice.ipAddress];
        selectDevicesForIntegrationMutation({
            variables: {
                integration: "kasa",
                devices: devices,
                addDevices: true,
                automationDevice: false,
            }
        })
        refetchData();
    }

    // Remove the selected variable from the interfaceDevices for kasa
    const removeSelectedInterface = (ipAddress) => {
        selectDevicesForIntegrationMutation({
            variables: {
                integration: "kasa",
                devices: [ipAddress],
                addDevices: false,
                automationDevice: false,
            }
        })
        refetchData();
    }

    // Add the selected variable to the interfaceDevices for kasa
    const addToSelectedAutomation = () => {
        devices = [selectedDevice.ipAddress];
        selectDevicesForIntegrationMutation({
            variables: {
                integration: "kasa",
                devices: devices,
                addDevices: true,
                automationDevice: true,
            }
        })
        refetchData();
    }

    // Remove the selected variable from the interfaceDevices for kasa
    const removeSelectedAutomation = (ipAddress) => {
        selectDevicesForIntegrationMutation({
            variables: {
                integration: "kasa",
                devices: [ipAddress],
                addDevices: false,
                automationDevice: true,
            }
        })
        refetchData();
    }

    const getNameForIp = (ipAddress) => {
        const device = [...integration.value.devices].filter(device => device.ipAddress === ipAddress);
        if (device !== undefined && device.length > 0 && device[0].name !== undefined) {
            return `${device[0].name}--${ipAddress}`;
        } else {
            return `${ipAddress}`;
        }
    }

    const clearSelected = () => {
        changeSelectedDevice("");
    }

    const updateSelected = (newDevice) => {
        changeSelectedDevice(newDevice);
    }

    var devices = [...integration.value.devices].sort((a, b) => (a.name > b.name) ? 1 : ((b.name > a.name) ? -1 : 0));
    var interfaceDevices = [...integration.value.interfaceDevices].map(device => getNameForIp(device)).sort((a, b) => (a > b) ? 1 : ((b > a) ? -1 : 0));
    var automatedDevices = [...integration.value.automatedDevices].map(device => getNameForIp(device)).sort((a, b) => (a > b) ? 1 : ((b > a) ? -1 : 0));

    return (
        <ExpansionPanel key="kasaExpansionPanel">
            <ExpansionPanelSummary key="kasa" >
                <div key="kasaTypographyNameDiv" className={classes.column}>
                    <Typography key="kasaTypographyName" className={classes.heading}>{integration.name}</Typography>
                </div>
                <div key="kasaTypographyHeadingDiv" className={classes.column}>
                    <Typography key="kasaTypographyHeading" className={classes.secondaryHeading}>Edit Configuration</Typography>
                </div>
            </ExpansionPanelSummary>
            <ExpansionPanelDetails key="kasaExpansionPanelDetails" className={classes.details}>
                <div key="kasaExpansionPanelDetailsDiv" className={classes.column}>
                    <FormControl key="kasaFormControl" className={classes.formControl}>
                        <InputLabel key="kasaFormControlInputLabel" id="demo-controlled-open-select-label">Kasa Devices</InputLabel>
                        <Select
                            key="kasaFormControlSelect"
                            id="demo-controlled-open-select"
                            value={selectedDevice}
                            defaultValue={selectedDevice}
                            onChange={handleChange}
                        >
                            <MenuItem key="kasaFormControlMenuItem" value="">
                                <em>None</em>
                            </MenuItem>
                            {
                                devices.map(device => <MenuItem key={"select_" + device.name} value={device}>{device.name}</MenuItem>)
                            }
                        </Select>
                    </FormControl>
                    {selectedDevice !== "" ?
                        <div className={classes.deviceDetails}>
                            <Typography
                                key={"selectedDeviceIpAddress" + selectedDevice.ipAddress}
                                className={classes.deviceHeading}
                                align="left">
                                <em className={classes.em}>ENTITY ID:</em>     {selectedDevice.ipAddress}
                            </Typography>
                            <Typography
                                key={"selectedDeviceType" + selectedDevice.ipAddress}
                                className={classes.deviceHeading}
                                align="left">
                                <em className={classes.em}>TYPE:</em>          {selectedDevice.type}
                            </Typography>
                            <Typography
                                key={"selectedDeviceAreaName" + selectedDevice.ipAddress}
                                className={classes.deviceHeading}
                                align="left">
                                <em className={classes.em}>AREA NAME:</em>     {selectedDevice.areaName}
                            </Typography>
                            <Typography
                                key={"selectedDeviceName" + selectedDevice.ipAddress}
                                className={classes.deviceHeading}
                                align="left">
                                <em className={classes.em}>NAME:</em>          {selectedDevice.name}
                            </Typography>
                            <div key={"selectedDeviceDiv" + selectedDevice.ipAddress} className={classes.buttonDiv}>
                                <Button key={"selectedDeviceInterfaceButton" + selectedDevice.ipAddress} className={classes.button} onClick={addToSelectedInterface}>Add to interface</Button>
                                <Button key={"selectedDeviceAutomatedButton" + selectedDevice.ipAddress} className={classes.button} onClick={addToSelectedAutomation}>Add to automated</Button>
                                <KasaEditMenu
                                    key={"selectedDeviceEditButton" + selectedDevice.ipAddress}
                                    device={selectedDevice}
                                    buttonStyle={classes.button}
                                    buttonText="edit device"
                                    areaData={areaData}
                                    refetchData={refetchData}
                                    updateSelected={updateSelected} />
                                <Button key={"selectedDeviceClearButton" + selectedDevice.ipAddress} className={classes.button} onClick={clearSelected}>clear</Button>
                            </div>
                        </div> : <div></div>}
                </div>
                <div key={"kasaInterfaceSelection" + selectedDevice.ipAddress} className={classNames(classes.column, classes.helper)}>
                    <Typography key={"kasaInterfaceSelectionTypography" + selectedDevice.ipAddress} className={classes.columnHeading}>Interface Devices</Typography>
                    {interfaceDevices.map(device =>
                        <Chip key={"interface_chip_" + device} label={device.split("--")[0]} className={classes.chip} onDelete={() => removeSelectedInterface(device.split("--")[1])} />
                    )}
                </div>
                <div key={"kasaAutomationSelection" + selectedDevice.ipAddress} className={classNames(classes.column, classes.helper)}>
                    <Typography key={"kasaAutomationSelectionTypography" + selectedDevice.ipAddress} className={classes.columnHeading}>Automation Devices</Typography>
                    {automatedDevices.map(device =>
                        <Chip key={"automation_chip_" + device} label={device.split("--")[0]} className={classes.chip} onDelete={() => removeSelectedAutomation(device.split("--")[1])} />
                    )}
                </div>
            </ExpansionPanelDetails>
        </ExpansionPanel>
    )
}
