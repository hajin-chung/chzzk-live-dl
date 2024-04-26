import { Dialog, DialogContent, DialogTrigger } from "./components/ui/dialog";
import { useState } from "react";
import { Button } from "./components/ui/button";
import { SettingsIcon } from "./components/icons";
import { Input } from "./components/ui/input";
import { updateCredentials } from "./lib/api";

export function Credentials() {
  const [ses, setSes] = useState("");
  const [aut, setAut] = useState("");

  function onSubmit() {
    updateCredentials(ses, aut);
  }

  return (
    <Dialog>
      <DialogTrigger>
        <Button variant="outline">
          <SettingsIcon />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <div className="flex flex-col gap-4 justify-end">
          <Input
            placeholder="NID_SES"
            value={ses}
            onChange={(e) => setSes(e.target.value)}
          />
          <Input
            placeholder="NID_AUT"
            value={aut}
            onChange={(e) => setAut(e.target.value)}
          />
        </div>
        <Button variant="outline" onClick={onSubmit}>
          Submit
        </Button>
      </DialogContent>
    </Dialog>
  );
}
