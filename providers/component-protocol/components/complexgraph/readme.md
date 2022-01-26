## Complex Graph Component

### Data Struct

|  Field   | Description  |
|  ----  | ----  |
| Title  | The title of line graph  |
| Dimensions  | Defining multiple legend in a Complex Graph |
| XAxis  | X-axis property definition |
| YAxis  | Y-axis property definition |
| Series  | Y-axis data |
| Inverse  | Inverse the x,y axis, the default is not inverse |

### Axis Struct

|  Field   | Description  |
|  ----  | ----  |
| Type  | Axis Type  |
| Name  | Axis Name  |
| Dimensions  | List of dimensions on the current axis  |
| Position  | Axis position,x-axis optional: bottom, top, y-axis optional left, right  |
| Data  | Indicates the x-axis data when the axis is the x-axis  |
| DataStructure  | Data structures on the same axis  |

### Sere Struct

|  Field   | Description  |
|  ----  | ----  |
| Name  | Sere Name  |
| Dimension  | data dimension  |
| Type  | Data display type, optional: bar,line  |
| Data  | Data |

### Demo

[Complex Graph Demo](../../examples/components/complexgraph_demo)