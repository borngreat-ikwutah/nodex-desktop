import { useState } from "react";
import { useNodeIdentity } from "./use-node-identity";
import {
  IconCopy,
  IconDiamondFilled,
  IconDeviceDesktop,
  IconClock,
  IconCircleFilled,
} from "@tabler/icons-react";
import { Card } from "@repo/ui/components/ui/card.tsx";
import { Button } from "@repo/ui/components/ui/button.tsx";
import { Input } from "@repo/ui/components/ui/input.tsx";
import { Label } from "@repo/ui/components/ui/label.tsx";
import {
  RadioGroup,
  RadioGroupItem,
} from "@repo/ui/components/ui/radio-group.tsx";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@repo/ui/components/ui/dialog.tsx";
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@repo/ui/components/ui/tabs.tsx";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@repo/ui/components/ui/table.tsx";

export function HomePage() {
  const { nodeID, isLoading } = useNodeIdentity();
  const [partnerId, setPartnerId] = useState("");
  const [mode, setMode] = useState<"screen" | "file">("screen");

  const displayID = nodeID || "000 000 000";

  return (
    <div className="flex h-full w-full flex-col font-sans overflow-y-auto select-none">
      <div className="flex flex-col items-center max-w-4xl w-full mx-auto p-4 md:p-8 space-y-8">
        {/* Top Section: Centralized ID */}
        <div className="flex flex-col items-center w-full space-y-6">
          <div className="space-y-3 flex flex-col items-center">
            <div className="bg-primary/10 w-12 h-12 rounded-xl flex items-center justify-center shadow-inner">
              <IconDiamondFilled className="size-6 text-primary" />
            </div>
            <div className="space-y-1 text-center">
              <h1 className="text-xl md:text-2xl font-bold tracking-tight text-foreground">
                Nodex Remote
              </h1>
              <p className="text-sm text-muted-foreground">
                Ready to receive secure connections.
              </p>
            </div>
          </div>

          {/* Centralized ID Card */}
          <Card className="w-full max-w-lg shadow-lg border-border/40 bg-card p-6 md:p-8 flex flex-col items-center space-y-8 rounded-[1.5rem]">
            <div className="space-y-3 w-full text-center">
              <Label className="text-muted-foreground text-xs uppercase tracking-widest font-semibold">
                Your ID
              </Label>
              <div className="flex items-center justify-center gap-4 group">
                <span className="text-3xl md:text-4xl leading-none font-bold tracking-tight text-foreground tabular-nums select-all cursor-text">
                  {isLoading ? "--- --- ---" : displayID}
                </span>
                <Button
                  variant="outline"
                  size="icon"
                  className="size-9 rounded-lg text-muted-foreground hover:text-foreground hover:bg-muted opacity-0 group-hover:opacity-100 transition-opacity"
                >
                  <IconCopy className="size-4" />
                </Button>
              </div>
            </div>

            <div className="w-full max-w-[200px] h-px bg-border/40" />

            <div className="space-y-3 w-full text-center">
              <Label className="text-muted-foreground text-xs uppercase tracking-widest font-semibold">
                Password
              </Label>
              <div className="flex items-center justify-center gap-4 group">
                <span className="text-3xl md:text-4xl leading-none font-bold tracking-[0.1em] text-foreground translate-y-1 select-all cursor-text">
                  ********
                </span>
                <Button
                  variant="outline"
                  size="icon"
                  className="size-9 rounded-lg text-muted-foreground hover:text-foreground hover:bg-muted opacity-0 group-hover:opacity-100 transition-opacity"
                >
                  <IconCopy className="size-4" />
                </Button>
              </div>
            </div>
          </Card>

          {/* Remote Connection Popup Trigger */}
          <Dialog>
            <DialogTrigger asChild>
              <Button className="h-11 px-8 text-sm rounded-xl font-medium shadow-md hover:scale-[1.02] active:scale-[0.98] transition-all">
                Connect to another computer
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[400px] rounded-[1.5rem] p-6">
              <DialogHeader className="space-y-1.5">
                <DialogTitle className="text-xl">Remote Connect</DialogTitle>
                <DialogDescription className="text-sm">
                  Enter the partner ID to start a secure session.
                </DialogDescription>
              </DialogHeader>
              <div className="py-4 space-y-6">
                <Input
                  type="text"
                  value={partnerId}
                  onChange={(e) => setPartnerId(e.target.value)}
                  placeholder="Partner ID"
                  className="w-full text-lg h-12 bg-muted/30 rounded-lg px-4 font-medium text-center tracking-wide placeholder:tracking-normal"
                />

                <RadioGroup
                  value={mode}
                  onValueChange={(val: any) => setMode(val)}
                  className="flex items-center justify-center gap-6"
                >
                  <div className="flex items-center space-x-2.5">
                    <RadioGroupItem value="screen" id="r1" className="size-4" />
                    <Label
                      htmlFor="r1"
                      className="cursor-pointer text-sm font-medium"
                    >
                      Screen control
                    </Label>
                  </div>
                  <div className="flex items-center space-x-2.5">
                    <RadioGroupItem value="file" id="r2" className="size-4" />
                    <Label
                      htmlFor="r2"
                      className="cursor-pointer text-sm text-muted-foreground hover:text-foreground font-medium transition-colors"
                    >
                      File transfer
                    </Label>
                  </div>
                </RadioGroup>
              </div>
              <Button className="w-full h-11 rounded-lg text-sm font-medium">
                Connect
              </Button>
            </DialogContent>
          </Dialog>
        </div>

        {/* Bottom Section: Tabs */}
        <div className="w-full max-w-2xl pt-4">
          <Tabs defaultValue="devices" className="w-full">
            <TabsList className="grid w-full grid-cols-2 max-w-[300px] mx-auto mb-6 bg-muted/40 p-1 rounded-lg h-auto">
              <TabsTrigger
                value="devices"
                className="rounded-md py-1.5 text-xs font-medium"
              >
                Recent Devices
              </TabsTrigger>
              <TabsTrigger
                value="sessions"
                className="rounded-md py-1.5 text-xs font-medium"
              >
                Recent Sessions
              </TabsTrigger>
            </TabsList>

            <TabsContent
              value="devices"
              className="animate-in fade-in-50 zoom-in-95 duration-200"
            >
              <Card className="rounded-xl border-border/40 shadow-sm overflow-hidden bg-card/40">
                <Table>
                  <TableHeader className="bg-muted/20">
                    <TableRow className="hover:bg-transparent">
                      <TableHead className="w-[200px] text-xs">
                        Device Name
                      </TableHead>
                      <TableHead className="text-xs">Status</TableHead>
                      <TableHead className="text-right text-xs">
                        Node ID
                      </TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    <TableRow>
                      <TableCell className="font-medium flex items-center gap-2.5 py-3 text-sm">
                        <div className="size-6 rounded-md bg-primary/10 flex items-center justify-center text-primary">
                          <IconDeviceDesktop className="size-3.5" />
                        </div>
                        Zhen's MacBook Pro
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                          <IconCircleFilled className="size-1.5 text-emerald-500" />
                          Online
                        </div>
                      </TableCell>
                      <TableCell className="text-right tabular-nums text-xs text-muted-foreground">
                        984 211 442
                      </TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell className="font-medium flex items-center gap-2.5 py-3 text-sm">
                        <div className="size-6 rounded-md bg-muted flex items-center justify-center text-muted-foreground">
                          <IconDeviceDesktop className="size-3.5" />
                        </div>
                        Office Workstation
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                          <IconCircleFilled className="size-1.5 text-muted-foreground/30" />
                          Offline
                        </div>
                      </TableCell>
                      <TableCell className="text-right tabular-nums text-xs text-muted-foreground">
                        221 094 112
                      </TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </Card>
            </TabsContent>

            <TabsContent
              value="sessions"
              className="animate-in fade-in-50 zoom-in-95 duration-200"
            >
              <Card className="rounded-xl border-border/40 shadow-sm overflow-hidden bg-card/40">
                <Table>
                  <TableHeader className="bg-muted/20">
                    <TableRow className="hover:bg-transparent">
                      <TableHead className="w-[200px] text-xs">
                        Connection
                      </TableHead>
                      <TableHead className="text-xs">Duration</TableHead>
                      <TableHead className="text-right text-xs">Date</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    <TableRow>
                      <TableCell className="font-medium flex items-center gap-2.5 py-3 text-sm">
                        <div className="size-6 rounded-md bg-blue-500/10 flex items-center justify-center text-blue-500">
                          <IconClock className="size-3.5" />
                        </div>
                        Remote to 984 211 442
                      </TableCell>
                      <TableCell className="text-xs text-muted-foreground">
                        45m 12s
                      </TableCell>
                      <TableCell className="text-right text-xs text-muted-foreground">
                        Today, 14:30
                      </TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell className="font-medium flex items-center gap-2.5 py-3 text-sm">
                        <div className="size-6 rounded-md bg-blue-500/10 flex items-center justify-center text-blue-500">
                          <IconClock className="size-3.5" />
                        </div>
                        Incoming from 221 094 112
                      </TableCell>
                      <TableCell className="text-xs text-muted-foreground">
                        1h 05m
                      </TableCell>
                      <TableCell className="text-right text-xs text-muted-foreground">
                        Yesterday
                      </TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </Card>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  );
}
