import { PacmanLoader } from "react-spinners";

const Loading = () => {
	return (
		<main className="flex justify-center items-center min-h-screen">
			<PacmanLoader color="#FFF200" size={75} />
		</main>
	);
}

export default Loading;