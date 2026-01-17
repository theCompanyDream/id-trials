import React from 'react';

import { Navigation } from '..';

const Layout = ({ children }) => (
	<article className='w-full h-full'>
		<Navigation />
		{children}
	</article>
)

export default Layout;