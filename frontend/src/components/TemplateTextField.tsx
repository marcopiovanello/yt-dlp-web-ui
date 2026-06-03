import { FC, useState } from 'react'

import DeleteIcon from '@mui/icons-material/Delete'
import EditIcon from '@mui/icons-material/Edit'
import {
  Button,
  Grid,
  TextField
} from '@mui/material'
import { useI18n } from '../hooks/useI18n'
import { CustomTemplate } from '../types'

interface Props {
  template: CustomTemplate
  onChange: (template: CustomTemplate) => void
  onDelete: (id: string) => void
}

const TemplateTextField: FC<Props> = ({ template, onChange, onDelete }) => {
  const { i18n } = useI18n()

  const [editedTemplate, setEditedTemplate] = useState(template)

  return (
    <Grid
      container
      spacing={2}
      justifyContent="center"
      alignItems="center"
      key={template.id}
      sx={{ mt: 1 }}
    >
      <Grid item xs={3}>
        <TextField
          fullWidth
          label={i18n.t('templatesEditorNameLabel')}
          defaultValue={template.name}
          onChange={(e) => setEditedTemplate({ ...editedTemplate, name: e.target.value })}
        />
      </Grid>
      <Grid item xs={9}>
        <TextField
          fullWidth
          label={i18n.t('templatesEditorContentLabel')}
          defaultValue={template.content}
          onChange={(e) => setEditedTemplate({ ...editedTemplate, content: e.target.value })}
          InputProps={{
            endAdornment: <div style={{ display: 'flex', gap: 2 }}>
              <Button
                variant='contained'
                onClick={() => onChange(editedTemplate)}>
                <EditIcon />
              </Button>
              <Button
                variant='contained'
                onClick={() => onDelete(editedTemplate.id)}>
                <DeleteIcon />
              </Button>
            </div>
          }}
        />
      </Grid>
    </Grid>
  )
}

export default TemplateTextField