import React from 'react';

import { Navigation } from '@backpack';

const Layout = ({ children }) => (
	<article className='w-full h-full'>
		<Navigation />
		{children}
	</article>
)

export default Layout;