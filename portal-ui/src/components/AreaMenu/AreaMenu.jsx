import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import { useMutation } from '@apollo/client';
import { updateIndexForArea } from './mutations/updateIndexForArea';

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

export default function AreaMenu({ area }) {
    const classes = useStyles();
    const [index, setIndex] = useState('');
    const [updateIndexMutation] = useMutation(updateIndexForArea);

    const [open, setOpen] = useState(false);

    const handleChange = ({ target }) => {
        updateIndexMutation({
            variables: {
                areaName: area.areaName,
                index: Number(target.value)
            }
        });
        handleClose();
    };


    if (area.index && area.index !== 999 && area.index !== index) {
        setIndex(area.index); 
    }

    const handleClickOpen = () => {
        setOpen(true);
    };

    const handleClose = () => {
        setOpen(false);
    };

    return (
        <div>
            <Button onClick={handleClickOpen}>Settings</Button>
            <Dialog open={open} onClose={handleClose}>
                <DialogTitle>{area.areaName} Settings</DialogTitle>
                <DialogContent>
                    <form className={classes.container} autoComplete='off'>
                        <FormControl className={classes.formControl}>
                            <InputLabel>Display Location</InputLabel>
                            <Select
                                id="priority-selection"
                                value={index}
                                onChange={handleChange}>
                                <MenuItem value={1}>1</MenuItem>
                                <MenuItem value={2}>2</MenuItem>
                                <MenuItem value={3}>3</MenuItem>
                                <MenuItem value={4}>4</MenuItem>
                                <MenuItem value={5}>5</MenuItem>
                                <MenuItem value={6}>6</MenuItem>
                                <MenuItem value={7}>7</MenuItem>
                                <MenuItem value={8}>8</MenuItem>
                                <MenuItem value={9}>9</MenuItem>
                            </Select>
                        </FormControl>
                    </form>
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
