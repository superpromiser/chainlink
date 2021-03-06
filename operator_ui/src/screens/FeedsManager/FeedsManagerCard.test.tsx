import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { FeedsManagerCard } from './FeedsManagerCard'
import { buildFeedsManagerFields } from 'support/factories/gql/fetchFeedsManagersWithProposals'
import { shortenHex } from 'src/utils/shortenHex'

const { getByRole, queryByText } = screen

function renderComponent(manager: FeedsManagerFields) {
  renderWithRouter(
    <>
      <Route path="/">
        <FeedsManagerCard manager={manager} />
      </Route>
      <Route path="/feeds_manager/edit">Redirect Success</Route>
    </>,
  )
}

describe('FeedsManagerCard', () => {
  it('renders a disconnected Feeds Manager', () => {
    const mgr = buildFeedsManagerFields({
      isBootstrapPeer: false,
      bootstrapPeerMultiaddr: '/dns4/blah',
    })

    renderComponent(mgr)

    expect(queryByText(mgr.name)).toBeInTheDocument()
    expect(queryByText(mgr.uri)).toBeInTheDocument()
    expect(queryByText(shortenHex(mgr.publicKey))).toBeInTheDocument()
    expect(queryByText('Flux Monitor')).toBeInTheDocument()
    expect(queryByText('Disconnected')).toBeInTheDocument()
    // We should not see the multiaddr because isBootstrapPeer is false
    expect(queryByText('/dns4/blah')).toBeNull()
  })

  it('renders a connected boostrapper Feeds Manager', () => {
    // Create a new manager with connected bootstrap values
    const mgr = buildFeedsManagerFields({
      jobTypes: [],
      isConnectionActive: true,
      isBootstrapPeer: true,
      bootstrapPeerMultiaddr: '/dns4/blah',
    })

    renderComponent(mgr)

    expect(queryByText(mgr.name)).toBeInTheDocument()
    expect(queryByText(mgr.uri)).toBeInTheDocument()
    expect(queryByText(shortenHex(mgr.publicKey))).toBeInTheDocument()
    expect(queryByText('Flux Monitor')).toBeNull()
    expect(queryByText('Connected')).toBeInTheDocument()
    expect(queryByText('/dns4/blah')).toBeInTheDocument()
  })

  it('navigates to edit', () => {
    renderComponent(buildFeedsManagerFields())

    userEvent.click(getByRole('link', { name: /edit/i }))

    expect(queryByText('Redirect Success')).toBeInTheDocument()
  })
})
