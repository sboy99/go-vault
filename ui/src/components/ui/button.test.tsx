import { describe, expect, it } from 'vitest'
import { render, screen } from '@solidjs/testing-library'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, THead, TBody, TR, TH, TD } from '@/components/ui/table'

describe('UI primitives', () => {
  it('renders Button', () => {
    render(() => <Button>Save</Button>)
    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument()
  })

  it('renders Badge', () => {
    render(() => <Badge color="success">success</Badge>)
    expect(screen.getByText('success')).toBeInTheDocument()
  })

  it('renders Table', () => {
    render(() => (
      <Table>
        <THead>
          <TR>
            <TH>Name</TH>
          </TR>
        </THead>
        <TBody>
          <TR>
            <TD>alpha.dump</TD>
          </TR>
        </TBody>
      </Table>
    ))
    expect(screen.getByText('alpha.dump')).toBeInTheDocument()
    expect(screen.getByText('Name')).toBeInTheDocument()
  })
})
