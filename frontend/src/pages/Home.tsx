import BasicTabs from '../components/Tabs';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import { Alert, Button, FilledInput } from '@mui/material';
import { useState } from 'react';
import { ShowFileServiceDialog } from '../../bindings/changeme/fileservice';

export default function HomePage() {

  const [showError, setShowError] = useState(false);
  const [error, setError] = useState("");
  const [selectedPath, setSelectedPath] = useState("");

  const onSelectFolderClick = async () => {
    try {
      const path = await ShowFileServiceDialog();
      setSelectedPath(path);
    } catch (e) {
      setShowError(true);
      setError("Cannot open folder dialog");
    }
  };

  return (
    <div>
      { showError && (
        <Alert severity="error">{error}</Alert>
      )}
      <FormControl sx={{ m: 1, width: '25ch' }} variant="filled">
                <InputLabel htmlFor={`folder-input`}>{selectedPath ? selectedPath : 'Select folder'}</InputLabel>
                <FilledInput
                  id={`folder-input`}
                  type='text'
        />
        <Button
          variant='outlined'
          onClick={onSelectFolderClick}
        >
          Select Folder
        </Button>
              </FormControl>
      <BasicTabs />
    </div>

  );
}
