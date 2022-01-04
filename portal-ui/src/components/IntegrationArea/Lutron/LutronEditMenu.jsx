import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import TextField from '@material-ui/core/TextField';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import { useMutation } from '@apollo/client';
import { editDeviceConfiguration } from '../mutations/editDeviceConfiguration';
import { useRef } from 'react';
// import { updateIndexForArea } from './mutations/updateIndexForArea';

const useStyles = makeStyles((theme) => ({
    container: {
        display: 'flex',
        flexWrap: 'wrap',
    },
    formControl: {
        margin: theme.spacing(1),
        minWidth: 200,
    },
}));

export default function LutronEditMenu({ device, buttonStyle, buttonText, areaData, refetchData, updateSelected }) {
    const classes = useStyles();
    const [areaName, setAreaName] = useState(device.areaName);
    const [updateDevice] = useMutation(editDeviceConfiguration);
    const nameRef = useRef(device.name);

    const [open, setOpen] = useState(false);

    const handleAreaNameChange = ({ target }) => {
        var newDevice = {
            areaName: target.value,
            name: device.name,
            type: device.type,
            lutronName: device.lutronName,
            id: device.id,
        };
        setAreaName(target.value);
        updateDevice({
            variables: {
                integration: "lutron",
                device: JSON.stringify(newDevice),
                removeEdit: false
            }
        });
        updateSelected(newDevice);
    };

    const handleNameChange = () => {
        var newDevice = {
            areaName: device.areaName,
            name: nameRef.current.value,
            type: device.type,
            lutronName: device.lutronName,
            id: device.id,
        };
        updateDevice({
            variables: {
                integration: "lutron",
                device: JSON.stringify(device),
                removeEdit: false
            }
        });
        updateSelected(newDevice);
    };

    const handleClickOpen = () => {
        setOpen(true);
    };

    const handleClose = () => {
        setOpen(false);
        refetchData();
    };

    return (
        <div>
            <Button onClick={handleClickOpen} className={buttonStyle}>{buttonText}</Button>
            <Dialog open={open} onClose={handleClose}>
                <DialogTitle>{device.entityId}</DialogTitle>
                <DialogContent>
                    <TextField
                        margin="dense"
                        id="ipAddress"
                        label="Unique ID"
                        defaultValue={device.id}
                        InputProps={{
                            readOnly: true,
                        }}
                        fullWidth
                    />
                    <TextField
                        margin="dense"
                        id="type"
                        label="Type"
                        defaultValue={device.type}
                        placeholder={device.type}
                        InputProps={{
                            readOnly: true,
                        }}
                        fullWidth
                    />
                    <FormControl
                        className={classes.formControl}
                        variant="filled"
                        fullWidth>
                        <InputLabel>Area Name</InputLabel>
                        <Select
                            id="priority-selection"
                            key={"selection_" + device.id}
                            value={areaName}
                            onChange={handleAreaNameChange}>
                            {areaData.areaNames.map(area => <MenuItem key={area + "lutronPopUp"} value={area}>{area}</MenuItem>)}
                        </Select>
                    </FormControl>
                    <TextField
                        variant="filled"
                        margin="dense"
                        id="name"
                        label="Device Name"
                        defaultValue={device.name}
                        placeholder={device.name}
                        onChange={handleNameChange}
                        helperText="Do you want a different name?"
                        inputRef={nameRef}
                        fullWidth
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose} color="primary">
                        Done
                    </Button>
                </DialogActions>
            </Dialog>
        </div>
    );
}