# About

This channelmodel (CM) library implements different IMT2020 channel models from ITU and 3GPP. You can use this model for simulation of wireless systems.

# Install
`go get github.com/wiless/channelmodel`


# Sample Usage
```
import "github.com/wiless/channelmodel"


	var rma CM.RMa
	plmodel:=CM.RMADefault()
  plmodel.SetFGz(0.7);
  rma.Set(plmodel)
  
  d=300; // in meters
  pl, err := rma.PLlos(dist)  // returns LOS-Pathloss for the distance 'd'
 
 
  pl, err = rma.PLnlos(dist)  // returns NLOS-Pathloss for the distance 'd'
 
 
 
  pl,islos,err = rma.PL(dist)  // returns Pathloss for the distance 'd' by randomly assigning LOS/NLOS, based on P_LOS distribution as in Specification 
 
 
   
```

# 
