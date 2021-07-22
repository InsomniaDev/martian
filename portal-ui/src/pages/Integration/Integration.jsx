import { makeStyles } from '@material-ui/core/styles';
import React from "react";
import { IntegrationArea } from '../../components/IntegrationArea/IntegrationArea';


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
    padding: `${theme.spacing.unit}px ${theme.spacing.unit * 2}px`,
  },
  link: {
    color: theme.palette.primary.main,
    textDecoration: 'none',
    '&:hover': {
      textDecoration: 'underline',
    },
  },
}));



export function Integration() {
//   const { loading, error, data, refetch,  } = useQuery(getMenuConfiguration, {
//     pollInterval: 500,
//     fetchPolicy: "no-cache"
//   });

  const classes = useStyles();

  // Set up the subscription for the light state
  // subscribeToMore({
  //   document: subscriptionForMenu,
  //   updateQuery: (prev, { subscriptionData: {
  //     data: { menuChange },
  //   } }) => {
  //     if (!menuChange) {
  //       return prev;
  //     }
  //     const newObject = {
  //       "menuConfiguration": menuChange
  //     }
  //     return Object.assign({}, prev, newObject)
  //   }
  // });

//   if (loading) return <p>loading ...</p>;
//   if (error) {
//     console.log(error);
//     return <p>Error :( {error}</p>;
//   }

  return (
    <div className={classes.root}>
      <IntegrationArea></IntegrationArea>
      <IntegrationArea></IntegrationArea>
      <IntegrationArea></IntegrationArea>
      <IntegrationArea></IntegrationArea>
    </div>
  )
}
