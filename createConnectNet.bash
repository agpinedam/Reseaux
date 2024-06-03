#!/bin/bash

# Definir las máquinas virtuales y sus interfaces
declare -A vms
vms=(
  [R1]="br1"
  [R2]="br1 br2 br3 br4"
  [R3]="br2 br5"
  [R4]="br3 br6"
  [R5]="br4 br5 br6 br7"
  [R6]="br7"
  [Client]="br1"
  [Server]="br7"
)

# Crear puentes y TAP interfaces
counter=1
for bridge in {1..7}; do
  BRIDGE_NAME="br$bridge"
  TAP_NAME="tap$counter"
  
  sudo ip tuntap add dev $TAP_NAME mode tap
  sudo ip link set $TAP_NAME up
  sudo brctl addbr $BRIDGE_NAME
  sudo brctl addif $BRIDGE_NAME $TAP_NAME
  sudo ip link set $BRIDGE_NAME up
  
  counter=$((counter + 1))
done

# Conectar las interfaces a las máquinas virtuales
for vm in "${!vms[@]}"; do
  IFS=' ' read -r -a bridges <<< "${vms[$vm]}"
  counter=1
  for bridge in "${bridges[@]}"; do
    INTERFACE_XML="/tmp/${vm}_interface_$counter.xml"
    cat <<EOF > $INTERFACE_XML
<interface type='bridge'>
  <mac address='52:54:00:$(dd if=/dev/urandom bs=1 count=3 2>/dev/null | hexdump -v -e '/1 ":%02X"')'/>
  <source bridge='$bridge'/>
  <model type='virtio'/>
  <address type='pci' domain='0x0000' bus='0x00' slot='0x03' function='0x0'/>
</interface>
EOF
    sudo virsh attach-device $vm $INTERFACE_XML --persistent
    rm $INTERFACE_XML
    counter=$((counter + 1))
  done
done

echo "Todas las interfaces de red virtuales han sido creadas y conectadas a las máquinas virtuales."
