import github from "../../assets/github.svg"
import { Navigation } from '@backpack';

const Layout = ({ children }) => (
	<article className='w-full h-full'>
		<Navigation
			links={[
				{link: "/analytics", name: "Analytics"},
				{link: "/about", name: "About"},
				{link: "/detail", name: "Detail"},
				{link: "https://user.tbrantleyii.dev/api/swagger/index.html", name: "Docs"},
				{link: "https://github.com/theCompanyDream/id-trials", img: github}
			]}
		/>
		{children}
	</article>
)

export default Layout;