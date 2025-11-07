import Link from "next/link"

const WhiteLinkButton = ({ where, what }) => {
    return (
        <Link
            href={where}
            className="link-secondary"
        >
            {what}
        </Link>
    )
}

export default WhiteLinkButton